package repos

import (
	"context"

	"github.com/inconshreveable/log15"

	"github.com/sourcegraph/sourcegraph/internal/api"
	dependenciesStore "github.com/sourcegraph/sourcegraph/internal/codeintel/dependencies/store"
	"github.com/sourcegraph/sourcegraph/internal/conf/reposource"
	"github.com/sourcegraph/sourcegraph/internal/database"
	"github.com/sourcegraph/sourcegraph/internal/database/dbutil"
	"github.com/sourcegraph/sourcegraph/internal/extsvc"
	"github.com/sourcegraph/sourcegraph/internal/extsvc/go"
	"github.com/sourcegraph/sourcegraph/internal/extsvc/go/gopackages"
	"github.com/sourcegraph/sourcegraph/internal/jsonc"
	"github.com/sourcegraph/sourcegraph/internal/types"
	"github.com/sourcegraph/sourcegraph/lib/errors"
	"github.com/sourcegraph/sourcegraph/schema"
)

// A GoPackagesSource creates git repositories from `*.tar.gz` files of
// published go dependencies from the Go ecosystem.
type GoPackagesSource struct {
	svc        *types.ExternalService
	connection schema.GoPackagesConnection
	depsStore  DependenciesStore
	client     npm.Client
}

// NewGoPackagesSource returns a new GoSource from the given external
// service.
func NewGoPackagesSource(svc *types.ExternalService) (*GoPackagesSource, error) {
	var c schema.GoPackagesConnection
	if err := jsonc.Unmarshal(svc.Config, &c); err != nil {
		return nil, errors.Errorf("external service id=%d config error: %s", svc.ID, err)
	}
	return &GoPackagesSource{
		svc:        svc,
		connection: c,
		/*dbStore initialized in SetDB */
		client: npm.NewHTTPClient(c.Registry, c.RateLimit, c.Credentials),
	}, nil
}

var _ Source = &GoPackagesSource{}

// ListRepos returns all go artifacts accessible to all connections
// configured in Sourcegraph via the external services configuration.
//
// [FIXME: deduplicate-listed-repos] The current implementation will return
// multiple repos with the same URL if there are different versions of it.
func (s *GoPackagesSource) ListRepos(ctx context.Context, results chan SourceResult) {
	goPackages, err := goPackages(s.connection)
	if err != nil {
		results <- SourceResult{Err: err}
		return
	}

	for _, goPackage := range goPackages {
		info, err := s.client.GetPackageInfo(ctx, goPackage)
		if err != nil {
			results <- SourceResult{Err: err}
			continue
		}

		repo := s.makeRepo(goPackage, info.Description)
		results <- SourceResult{
			Source: s,
			Repo:   repo,
		}
	}

	totalDBFetched, totalDBResolved, lastID := 0, 0, 0
	pkgVersions := map[string]*gopkg.PackageInfo{}
	for {
		dbDeps, err := s.depsStore.ListDependencyRepos(ctx, dependenciesStore.ListDependencyReposOpts{
			Scheme:      dependenciesStore.GoPackagesScheme,
			After:       lastID,
			Limit:       100,
			NewestFirst: true,
		})
		if err != nil {
			results <- SourceResult{Err: err}
			return
		}
		if len(dbDeps) == 0 {
			break
		}
		totalDBFetched += len(dbDeps)
		lastID = dbDeps[len(dbDeps)-1].ID
		for _, dbDep := range dbDeps {
			parsedDbPackage, err := reposource.ParseGoPackageFromPackageSyntax(dbDep.Name)
			if err != nil {
				log15.Error("failed to parse go package name retrieved from database", "package", dbDep.Name, "error", err)
				continue
			}

			goDependency := reposource.GoDependency{GoPackage: parsedDbPackage, Version: dbDep.Version}
			pkgKey := goDependency.PackageSyntax()
			info := pkgVersions[pkgKey]

			if info == nil {
				info, err = s.client.GetPackageInfo(ctx, goDependency.GoPackage)
				if err != nil {
					pkgVersions[pkgKey] = &gopkg.PackageInfo{Versions: map[string]*gopkg.DependencyInfo{}}
					continue
				}

				pkgVersions[pkgKey] = info
			}

			if _, hasVersion := info.Versions[goDependency.Version]; !hasVersion {
				continue
			}

			repo := s.makeRepo(goDependency.GoPackage, info.Description)
			totalDBResolved++
			results <- SourceResult{Source: s, Repo: repo}
		}
	}
	log15.Info("finish resolving go artifacts", "totalDB", totalDBFetched, "totalDBResolved", totalDBResolved, "totalConfig", len(goPackages))
}

func (s *GoPackagesSource) GetRepo(ctx context.Context, name string) (*types.Repo, error) {
	pkg, err := reposource.ParseGoPackageFromRepoURL(name)
	if err != nil {
		return nil, err
	}

	info, err := s.client.GetPackageInfo(ctx, pkg)
	if err != nil {
		return nil, err
	}

	return s.makeRepo(pkg, info.Description), nil
}

func (s *GoPackagesSource) makeRepo(goPackage *reposource.GoPackage, description string) *types.Repo {
	urn := s.svc.URN()
	cloneURL := goPackage.CloneURL()
	repoName := goPackage.RepoName()
	return &types.Repo{
		Name:        repoName,
		Description: description,
		URI:         string(repoName),
		ExternalRepo: api.ExternalRepoSpec{
			ID:          string(repoName),
			ServiceID:   extsvc.TypeGoPackages,
			ServiceType: extsvc.TypeGoPackages,
		},
		Private: false,
		Sources: map[string]*types.SourceInfo{
			urn: {
				ID:       urn,
				CloneURL: cloneURL,
			},
		},
		Metadata: &gopackages.Metadata{
			Package: goPackage,
		},
	}
}

// ExternalServices returns a singleton slice containing the external service.
func (s *GoPackagesSource) ExternalServices() types.ExternalServices {
	return types.ExternalServices{s.svc}
}

func (s *GoPackagesSource) SetDB(db dbutil.DB) {
	s.depsStore = dependenciesStore.GetStore(database.NewDB(db))
}

// goPackages gets the list of applicable packages by de-duplicating dependencies
// present in the configuration.
func goPackages(connection schema.GoPackagesConnection) ([]*reposource.GoPackage, error) {
	dependencies, err := goDependencies(connection)
	if err != nil {
		return nil, err
	}
	goPackages := []*reposource.GoPackage{}
	isAdded := make(map[string]bool)
	for _, dep := range dependencies {
		if key := dep.PackageSyntax(); !isAdded[key] {
			goPackages = append(goPackages, dep.GoPackage)
			isAdded[key] = true
		}
	}
	return goPackages, nil
}

func goDependencies(connection schema.GoPackagesConnection) (dependencies []*reposource.GoDependency, err error) {
	for _, dep := range connection.Dependencies {
		dependency, err := reposource.ParseGoDependency(dep)
		if err != nil {
			return nil, err
		}
		dependencies = append(dependencies, dependency)
	}
	return dependencies, nil
}

type DependenciesStore interface {
	ListDependencyRepos(ctx context.Context, opts dependenciesStore.ListDependencyReposOpts) ([]dependenciesStore.DependencyRepo, error)
}
