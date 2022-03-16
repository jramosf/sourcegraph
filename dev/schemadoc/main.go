package main

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/sync/errgroup"

	connections "github.com/sourcegraph/sourcegraph/internal/database/connections/live"
	migrationstore "github.com/sourcegraph/sourcegraph/internal/database/migration/store"
	"github.com/sourcegraph/sourcegraph/internal/lazyregexp"
	"github.com/sourcegraph/sourcegraph/internal/observation"
	"github.com/sourcegraph/sourcegraph/lib/errors"
)

type runFunc func(quiet bool, cmd ...string) (string, error)

const databaseNamePrefix = "schemadoc-gen-temp-"

const containerName = "schemadoc"

var logger = log.New(os.Stderr, "", log.LstdFlags)

var versionRe = lazyregexp.New(`\b12\.\d+\b`)

type databaseFactory func(dsn string, appName string, observationContext *observation.Context) (*sql.DB, error)

var schemas = map[string]struct {
	destinationFilename string
	factory             databaseFactory
}{
	"frontend":  {"schema.md", connections.MigrateNewFrontendDB},
	"codeintel": {"schema.codeintel.md", connections.MigrateNewCodeIntelDB},
}

// This script generates markdown formatted output containing descriptions of
// the current dabase schema, obtained from postgres. The correct PGHOST,
// PGPORT, PGUSER etc. env variables must be set to run this script.
func main() {
	if err := mainErr(); err != nil {
		log.Fatal(err)
	}
}

func mainErr() error {
	// Run pg12 locally if it exists
	if version, _ := exec.Command("psql", "--version").CombinedOutput(); versionRe.Match(version) {
		return mainLocal()
	}

	return mainContainer()
}

func mainLocal() error {
	dataSourcePrefix := "dbname=" + databaseNamePrefix

	g, _ := errgroup.WithContext(context.Background())
	for name, schema := range schemas {
		name, schema := name, schema
		g.Go(func() error {
			return generateAndWrite(name, schema.factory, dataSourcePrefix+name, nil, schema.destinationFilename)
		})
	}

	return g.Wait()
}

func mainContainer() error {
	logger.Printf("Running PostgreSQL 12 in docker")

	prefix, shutdown, err := startDocker()
	if err != nil {
		log.Fatal(err)
	}
	defer shutdown()

	dataSourcePrefix := "postgres://postgres@127.0.0.1:5433/postgres?dbname=" + databaseNamePrefix

	g, _ := errgroup.WithContext(context.Background())
	for name, schema := range schemas {
		name, schema := name, schema
		g.Go(func() error {
			return generateAndWrite(name, schema.factory, dataSourcePrefix+name, prefix, schema.destinationFilename)
		})
	}

	return g.Wait()
}

func generateAndWrite(name string, factory databaseFactory, dataSource string, commandPrefix []string, destinationFile string) error {
	run := runWithPrefix(commandPrefix)

	// Try to drop a database if it already exists
	_, _ = run(true, "dropdb", databaseNamePrefix+name)

	// Let's also try to clean up after ourselves
	defer func() { _, _ = run(true, "dropdb", databaseNamePrefix+name) }()

	if out, err := run(false, "createdb", databaseNamePrefix+name); err != nil {
		return errors.Wrap(err, fmt.Sprintf("run: %s", out))
	}

	db, err := factory(dataSource, "", &observation.TestContext)
	if err != nil {
		return err
	}
	defer db.Close()

	out, err := generateInternal(db, name, run)
	if err != nil {
		return err
	}

	return os.WriteFile(destinationFile, []byte(out), os.ModePerm)
}

func startDocker() (commandPrefix []string, shutdown func(), _ error) {
	if err := exec.Command("docker", "image", "inspect", "postgres:12").Run(); err != nil {
		logger.Println("docker pull postgres:12")
		pull := exec.Command("docker", "pull", "postgres:12")
		pull.Stdout = logger.Writer()
		pull.Stderr = logger.Writer()
		if err := pull.Run(); err != nil {
			return nil, nil, errors.Wrap(err, "docker pull postgres:12")
		}
		logger.Println("docker pull complete")
	}

	run := runWithPrefix(nil)

	_, _ = run(true, "docker", "rm", "--force", containerName)
	server := exec.Command("docker", "run", "--rm", "--name", containerName, "-e", "POSTGRES_HOST_AUTH_METHOD=trust", "-p", "5433:5432", "postgres:12")
	if err := server.Start(); err != nil {
		return nil, nil, errors.Wrap(err, "docker run")
	}

	shutdown = func() {
		_ = server.Process.Kill()
		_, _ = run(true, "docker", "kill", containerName)
		_ = server.Wait()
	}

	attempts := 0
	for {
		attempts++
		// TODO - not sure why this would work...?
		if err := exec.Command("pg_isready", "-U", "postgres", "-h", "127.0.0.1", "-p", "5433").Run(); err == nil {
			break
		} else if attempts > 30 {
			shutdown()
			return nil, nil, errors.Wrap(err, "pg_isready timeout")
		}
		time.Sleep(time.Second)
	}

	return []string{"docker", "exec", "-u", "postgres", containerName}, shutdown, nil
}

func generateInternal(db *sql.DB, name string, run runFunc) (_ string, err error) {
	store := migrationstore.NewWithDB(db, "schema_migrations", migrationstore.NewOperations(&observation.TestContext))
	schemas, err := store.Describe(context.Background())
	if err != nil {
		return "", err
	}

	docs := []string{}
	types := map[string][]string{}

	for schemaName, schema := range schemas {
		sort.Slice(schema.Tables, func(i, j int) bool { return schema.Tables[i].Name < schema.Tables[j].Name })
		for _, table := range schema.Tables {
			docs = append(docs, fmt.Sprintf("# Table \"%s.%s\"", schemaName, table.Name))
			docs = append(docs, "\n")
			if table.Comment != "" {
				docs = append(docs, table.Comment+"\n")
			}

			docs = append(docs, "| "+strings.Join([]string{"Column", "Type", "Nullable", "Default", "Comment"}, " | ")+" |")
			docs = append(docs, "| "+strings.Join([]string{"---", "---", "---", "---", "---"}, " | ")+" |")

			sort.Slice(table.Columns, func(i, j int) bool { return table.Columns[i].Name < table.Columns[j].Name })
			for _, column := range table.Columns {
				values := []string{}
				values = append(values, column.Name)
				values = append(values, column.TypeName)
				if column.IsNullable {
					values = append(values, "Yes")
				} else {
					values = append(values, "No")
				}
				if column.IsGenerated == "ALWAYS" {
					// TODO - parse more specifically
					values = append(values, "generated always as ("+column.GenerationExpression+") stored")
				} else {
					values = append(values, column.Default)
				}

				values = append(values, column.Comment)
				docs = append(docs, "| "+strings.Join(values, " | ")+" |")
			}

			if len(table.Indexes) > 0 {
				docs = append(docs, "### Indexes")
				docs = append(docs, "| "+strings.Join([]string{"Name", "IsPrimaryKey", "IsUnique", "IsExclusion", "IsDeferrable", "IndexDefinition"}, " | ")+" |")
				docs = append(docs, "| "+strings.Join([]string{"---", "---", "---", "---", "---", "---"}, " | ")+" |")
			}
			sort.Slice(table.Indexes, func(i, j int) bool {
				return table.Indexes[i].Name < table.Indexes[j].Name
			})
			for _, index := range table.Indexes {
				values := []string{}
				values = append(values, index.Name)
				if index.IsPrimaryKey {
					values = append(values, "Yes")
				} else {
					values = append(values, "no")
				}
				if index.IsUnique {
					values = append(values, "Yes")
				} else {
					values = append(values, "no")
				}
				if index.IsExclusion {
					values = append(values, "Yes")
				} else {
					values = append(values, "no")
				}
				if index.IsDeferrable {
					values = append(values, "Yes")
				} else {
					values = append(values, "no")
				}
				values = append(values, index.IndexDefinition)

				docs = append(docs, "| "+strings.Join(values, " | ")+" |")
			}

			numCheckConstraints := 0
			numForeignKeyConstraints := 0
			for _, constraint := range table.Constraints {
				switch constraint.ConstraintType {
				case "c":
					numCheckConstraints++
				case "f":
					numForeignKeyConstraints++
				}
			}

			if numCheckConstraints > 0 {
				docs = append(docs, "### Check constraints")
				docs = append(docs, "| "+strings.Join([]string{"Name", "Definition"}, " | ")+" |")
				docs = append(docs, "| "+strings.Join([]string{"---", "---"}, " | ")+" |")
			}
			for _, constraint := range table.Constraints {
				if constraint.ConstraintType == "c" {
					docs = append(docs, "| "+strings.Join([]string{
						constraint.Name,
						constraint.ConstraintDefinition,
					}, " | ")+" |")
				}
			}
			if numForeignKeyConstraints > 0 {
				docs = append(docs, "### Foreign key constraints")
				docs = append(docs, "| "+strings.Join([]string{"Name", "References", "Definition"}, " | ")+" |")
				docs = append(docs, "| "+strings.Join([]string{"---", "---", "---"}, " | ")+" |")
			}
			for _, constraint := range table.Constraints {
				if constraint.ConstraintType == "f" {
					docs = append(docs, "| "+strings.Join([]string{
						constraint.Name,
						constraint.RefTableName, // TODO - make a link
						constraint.ConstraintDefinition,
					}, " | ")+" |")
				}
			}

			if len(table.Triggers) > 0 {
				docs = append(docs, "### Triggers")
				docs = append(docs, "| "+strings.Join([]string{"Name", "Definition"}, " | ")+" |")
				docs = append(docs, "| "+strings.Join([]string{"---", "---"}, " | ")+" |")
			}
			for _, trigger := range table.Triggers {
				docs = append(docs, "| "+strings.Join([]string{
					trigger.Name,
					trigger.Definition,
				}, " | ")+" |")
			}

			type tableAndConstraint struct {
				migrationstore.Table
				migrationstore.Constraint
			}
			tableAndConstraints := []tableAndConstraint{}

			for _, otherTable := range schema.Tables {
				for _, constraint := range otherTable.Constraints {
					if constraint.RefTableName == table.Name {
						tableAndConstraints = append(tableAndConstraints, tableAndConstraint{otherTable, constraint})
					}
				}
			}
			sort.Slice(tableAndConstraints, func(i, j int) bool {
				return tableAndConstraints[i].Constraint.Name < tableAndConstraints[j].Constraint.Name
			})
			if len(tableAndConstraints) > 0 {
				docs = append(docs, "### References")
				docs = append(docs, "| "+strings.Join([]string{"Name", "Definition"}, " | ")+" |")
				docs = append(docs, "| "+strings.Join([]string{"---", "---"}, " | ")+" |")
			}
			for _, tableAndConstraint := range tableAndConstraints {
				docs = append(docs, "| "+strings.Join([]string{
					tableAndConstraint.Table.Name, // TOOD - make link
					tableAndConstraint.Constraint.Name,
				}, " | ")+" |")
			}
		}

		sort.Slice(schema.Views, func(i, j int) bool { return schema.Views[i].Name < schema.Views[j].Name })
		for _, view := range schema.Views {
			docs = append(docs, fmt.Sprintf("# View \"public.%s\"\n", view.Name))
			docs = append(docs, fmt.Sprintf("## View query:\n\n```sql\n%s\n```\n", view.Definition))
		}

		for _, enum := range schema.Enums {
			types[enum.Name] = enum.Labels
		}
	}

	combined := strings.Join(docs, "\n")

	if len(types) > 0 {
		buf := bytes.NewBuffer(nil)
		buf.WriteString("\n")

		var typeNames []string
		for k := range types {
			typeNames = append(typeNames, k)
		}
		sort.Strings(typeNames)

		for _, name := range typeNames {
			buf.WriteString("# Type ")
			buf.WriteString(name)
			buf.WriteString("\n\n- ")
			buf.WriteString(strings.Join(types[name], "\n- "))
			buf.WriteString("\n\n")
		}

		combined += buf.String()
	}

	return combined, nil
}

func runWithPrefix(prefix []string) runFunc {
	return func(quiet bool, cmd ...string) (string, error) {
		cmd = append(prefix, cmd...)

		c := exec.Command(cmd[0], cmd[1:]...)
		if !quiet {
			c.Stderr = logger.Writer()
		}

		out, err := c.Output()
		return string(out), err
	}
}
