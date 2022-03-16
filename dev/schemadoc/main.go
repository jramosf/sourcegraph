package main

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"sort"
	"strconv"
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
			sizes := []int{
				len("Column"),
				len("Type"),
				len("Collation"),
				len("Nullable"),
				len("Default"),
			}
			for _, column := range table.Columns {
				if n := len(column.Name); n > sizes[0] {
					sizes[0] = n
				}
				if n := len(column.TypeName); n > sizes[1] {
					sizes[1] = n
				}
				defaultValue := column.Default
				if column.IsGenerated == "ALWAYS" {
					defaultValue = "generated always as (" + column.GenerationExpression + ") stored"
				}
				if n := len(defaultValue); n > sizes[4] {
					sizes[4] = n
				}
			}

			center := func(s string, n int) string {
				x := float64(n - len(s))
				i := int(math.Floor(x / 2))
				if i <= 0 {
					i = 1
				}
				j := int(math.Ceil(x / 2))
				if j <= 0 {
					j = 1
				}

				return strings.Repeat(" ", i) + s + strings.Repeat(" ", j)
			}

			header := strings.Join([]string{
				center("Column", sizes[0]+2),
				center("Type", sizes[1]+2),
				center("Collation", sizes[2]+2),
				center("Nullable", sizes[3]+2),
				center("Default", sizes[4]+2),
			}, "|")

			sep := strings.Join([]string{
				strings.Repeat("-", sizes[0]+2),
				strings.Repeat("-", sizes[1]+2),
				strings.Repeat("-", sizes[2]+2),
				strings.Repeat("-", sizes[3]+2),
				strings.Repeat("-", sizes[4]+2),
			}, "+")

			docs = append(docs, fmt.Sprintf("# Table \"%s.%s\"", schemaName, table.Name))
			docs = append(docs, "```")
			docs = append(docs, header)
			docs = append(docs, sep)

			sort.Slice(table.Columns, func(i, j int) bool { return table.Columns[i].Index < table.Columns[j].Index })

			for _, column := range table.Columns {
				nullConstraint := "not null"
				if column.IsNullable {
					nullConstraint = ""
				}

				defaultValue := column.Default
				if column.IsGenerated == "ALWAYS" {
					defaultValue = "generated always as (" + column.GenerationExpression + ") stored"
				}

				col := " " + strings.Join([]string{
					fmt.Sprintf("%-"+strconv.Itoa(sizes[0])+"s", column.Name),
					fmt.Sprintf("%-"+strconv.Itoa(sizes[1])+"s", column.TypeName),
					fmt.Sprintf("%-"+strconv.Itoa(sizes[2])+"s", ""),
					fmt.Sprintf("%-"+strconv.Itoa(sizes[3])+"s", nullConstraint),
					defaultValue,
				}, " | ")

				docs = append(docs, col)
			}

			if len(table.Indexes) > 0 {
				docs = append(docs, "Indexes:")
			}
			sort.Slice(table.Indexes, func(i, j int) bool {
				if table.Indexes[i].IsUnique && !table.Indexes[j].IsUnique {
					return true
				}
				if !table.Indexes[i].IsUnique && table.Indexes[j].IsUnique {
					return false
				}
				return table.Indexes[i].Name < table.Indexes[j].Name
			})
			for _, index := range table.Indexes {
				if !index.IsPrimaryKey {
					continue
				}

				deferrable := ""
				if index.IsDeferrable {
					// deferrable = " DEFERRABLE"
				}
				def := strings.TrimSpace(strings.Split(index.IndexDefinition, "USING")[1])
				docs = append(docs, fmt.Sprintf("    %q PRIMARY KEY, %s%s", index.Name, def, deferrable))
			}
			for _, index := range table.Indexes {
				if index.IsPrimaryKey {
					continue
				}

				uq := ""
				if index.IsUnique {
					uq = " UNIQUE CONSTRAINT,"
				}
				deferrable := ""
				if index.IsDeferrable {
					deferrable = " DEFERRABLE"
				}
				def := strings.TrimSpace(strings.Split(index.IndexDefinition, "USING")[1])
				if index.IsExclusion {
					def = "EXCLUDE USING " + def
				}
				docs = append(docs, fmt.Sprintf("    %q%s %s%s", index.Name, uq, def, deferrable))
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
				docs = append(docs, "Check constraints:")
			}
			for _, constraint := range table.Constraints {
				if constraint.ConstraintType == "c" {
					deferrable := ""
					if constraint.IsDeferrable {
						// deferrable = " DEFERRABLE"
					}
					docs = append(docs, fmt.Sprintf("    %q %s%s", constraint.Name, constraint.ConstraintDefinition, deferrable))
				}

			}
			if numForeignKeyConstraints > 0 {
				docs = append(docs, "Foreign-key constraints:")
			}
			for _, constraint := range table.Constraints {
				if constraint.ConstraintType == "f" {
					deferrable := ""
					if constraint.IsDeferrable {
						// deferrable = " DEFERRABLE"
					}
					docs = append(docs, fmt.Sprintf("    %q %s%s", constraint.Name, constraint.ConstraintDefinition, deferrable))
				}
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
				// if tableAndConstraints[i].Table.Name == tableAndConstraints[j].Table.Name {
				return tableAndConstraints[i].Constraint.Name < tableAndConstraints[j].Constraint.Name
				// }
				// return tableAndConstraints[i].Table.Name < tableAndConstraints[j].Table.Name
			})
			if len(tableAndConstraints) > 0 {
				docs = append(docs, "Referenced by:")
			}
			for _, tableAndConstraint := range tableAndConstraints {
				docs = append(docs, fmt.Sprintf("    TABLE %q CONSTRAINT %q %s", tableAndConstraint.Table.Name, tableAndConstraint.Constraint.Name, tableAndConstraint.Constraint.ConstraintDefinition))
			}

			if len(table.Triggers) > 0 {
				docs = append(docs, "Triggers:")
			}
			for _, trigger := range table.Triggers {
				def := strings.TrimSpace(strings.SplitN(trigger.Definition, trigger.Name, 2)[1])
				docs = append(docs, fmt.Sprintf("    %s %s", trigger.Name, def))
			}

			docs = append(docs, "\n```\n")

			if table.Comment != "" {
				docs = append(docs, table.Comment+"\n")
			}

			sort.Slice(table.Columns, func(i, j int) bool { return table.Columns[i].Name < table.Columns[j].Name })
			for _, column := range table.Columns {
				if column.Comment != "" {
					docs = append(docs, fmt.Sprintf("**%s**: %s\n", column.Name, column.Comment))
				}
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
