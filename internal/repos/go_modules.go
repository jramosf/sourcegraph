package repos

import (
	"context"
	"github.com/inconshreveable/log15"
	"github.com/sourcegraph/sourcegraph/internal/extsvc/gomodproxy"
	"github.com/sourcegraph/sourcegraph/internal/httpcli"

	"github.com/sourcegraph/sourcegraph/internal/api"
	dependenciesStore "github.com/sourcegraph/sourcegraph/internal/codeintel/dependencies/store"
	"github.com/sourcegraph/sourcegraph/internal/conf/reposource"
	"github.com/sourcegraph/sourcegraph/internal/database"
	"github.com/sourcegraph/sourcegraph/internal/database/dbutil"
	"github.com/sourcegraph/sourcegraph/internal/extsvc"
	"github.com/sourcegraph/sourcegraph/internal/extsvc/go/gopackages"
	"github.com/sourcegraph/sourcegraph/internal/jsonc"
	"github.com/sourcegraph/sourcegraph/internal/types"
	"github.com/sourcegraph/sourcegraph/lib/errors"
	"github.com/sourcegraph/sourcegraph/schema"
)

// A GoModulesSource creates git repositories from go module zip files of
// published go dependencies from the Go ecosystem.
type GoModulesSource struct {
	svc       *types.ExternalService
	config    schema.GoModuleProxiesConnection
	depsStore DependenciesStore
	client    *gomodproxy.Client
}

// NewGoModulesSource returns a new GoModulesSource from the given external service.
func NewGoModulesSource(svc *types.ExternalService, cf *httpcli.Factory) (*GoModulesSource, error) {
	var c schema.GoModuleProxiesConnection
	if err := jsonc.Unmarshal(svc.Config, &c); err != nil {
		return nil, errors.Errorf("external service id=%d config error: %s", svc.ID, err)
	}

	cli, err := cf.Doer()
	if err != nil {
		return nil, err
	}

	return &GoModulesSource{
		svc:    svc,
		config: c,
		/*dbStore initialized in SetDB */
		client: gomodproxy.NewClient(&c, cli),
	}, nil
}

var _ Source = &GoModulesSource{}

func (s *GoModulesSource) ListRepos(ctx context.Context, results chan SourceResult) {
	goModules, err := goModules(s.connection)
	if err != nil {
		results <- SourceResult{Err: err}
		return
	}

	for _, goPackage := range goModules {
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
			Scheme:      dependenciesStore.GoModulesScheme,
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
			parsedDbPackage, err := reposource.ParseGoModuleFromPackageSyntax(dbDep.Name)
			if err != nil {
				log15.Error("failed to parse go package name retrieved from database", "package", dbDep.Name, "error", err)
				continue
			}

			goDependency := reposource.GoDependency{GoModule: parsedDbPackage, Version: dbDep.Version}
			pkgKey := goDependency.PackageSyntax()
			info := pkgVersions[pkgKey]

			if info == nil {
				info, err = s.client.GetPackageInfo(ctx, goDependency.GoModule)
				if err != nil {
					pkgVersions[pkgKey] = &gopkg.PackageInfo{Versions: map[string]*gopkg.DependencyInfo{}}
					continue
				}

				pkgVersions[pkgKey] = info
			}

			if _, hasVersion := info.Versions[goDependency.Version]; !hasVersion {
				continue
			}

			repo := s.makeRepo(goDependency.GoModule, info.Description)
			totalDBResolved++
			results <- SourceResult{Source: s, Repo: repo}
		}
	}
	log15.Info("finish resolving go artifacts", "totalDB", totalDBFetched, "totalDBResolved", totalDBResolved, "totalConfig", len(goModules))
}

func (s *GoModulesSource) GetRepo(ctx context.Context, name string) (*types.Repo, error) {
	pkg, err := reposource.ParseGoModuleFromRepoURL(name)
	if err != nil {
		return nil, err
	}

	info, err := s.client.GetPackageInfo(ctx, pkg)
	if err != nil {
		return nil, err
	}

	return s.makeRepo(pkg, info.Description), nil
}

func (s *GoModulesSource) makeRepo(goPackage *reposource.GoModule, description string) *types.Repo {
	urn := s.svc.URN()
	cloneURL := goPackage.CloneURL()
	repoName := goPackage.RepoName()
	return &types.Repo{
		Name:        repoName,
		Description: description,
		URI:         string(repoName),
		ExternalRepo: api.ExternalRepoSpec{
			ID:          string(repoName),
			ServiceID:   extsvc.TypeGoModules,
			ServiceType: extsvc.TypeGoModules,
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
func (s *GoModulesSource) ExternalServices() types.ExternalServices {
	return types.ExternalServices{s.svc}
}

func (s *GoModulesSource) SetDB(db dbutil.DB) {
	s.depsStore = dependenciesStore.GetStore(database.NewDB(db))
}

// goModules gets the list of modules by de-duplicating dependencies
func goModules(connection schema.GoModulesConnection) ([]*reposource.GoModule, error) {
	dependencies, err := goDependencies(connection)
	if err != nil {
		return nil, err
	}
	goModules := []*reposource.GoModule{}
	isAdded := make(map[string]bool)
	for _, dep := range dependencies {
		if key := dep.PackageSyntax(); !isAdded[key] {
			goModules = append(goModules, dep.GoModule)
			isAdded[key] = true
		}
	}
	return goModules, nil
}

func goDependencies(connection schema.GoModulesConnection) (dependencies []*reposource.GoDependency, err error) {
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
