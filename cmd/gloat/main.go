package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/gsamokovarov/gloat"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

const usage = `Usage gloat: [OPTION ...] [COMMAND ...]

Gloat is a Go SQL migration utility.

Commands:
  new           Create a new migration folder
  up            Apply new migrations
  down          Revert the last applied migration

Options:
  -src          The folder with migrations
                (default $DATABASE_SRC" or db/migrations)
  -url          The database connection URL
                (default $DATABASE_URL)
  -help         Show this message
`

type arguments struct {
	url  string
	src  string
	rest []string
}

func main() {
	args := parseArguments()

	var cmdName string
	if len(args.rest) > 0 {
		cmdName = args.rest[0]
	}

	var err error
	switch cmdName {
	case "up":
		err = upCmd(args)
	case "down":
		err = downCmd(args)
	case "new":
		err = newCmd(args)
	default:
		fmt.Fprintf(os.Stderr, usage)
		os.Exit(2)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %+v\n", err)
		os.Exit(2)
	}
}

func upCmd(args arguments) error {
	gl, err := setupGloat(args)
	if err != nil {
		return err
	}

	migrations, err := gl.Unapplied()
	if err != nil {
		return err
	}

	appliedMigrations := map[int64]bool{}

	for _, migration := range migrations {
		fmt.Printf("Applying: %d...\n", migration.Version)

		if err := gl.Apply(migration); err != nil {
			return err
		}

		appliedMigrations[migration.Version] = true
	}

	if len(appliedMigrations) == 0 {
		fmt.Printf("No migrations to apply\n")
	}

	return nil
}

func downCmd(args arguments) error {
	gl, err := setupGloat(args)
	if err != nil {
		return err
	}

	migration, err := gl.Current()
	if err != nil {
		return err
	}

	if migration == nil {
		fmt.Printf("No migrations to apply\n")
		return nil
	}

	fmt.Printf("Reverting: %d...\n", migration.Version)

	if err := gl.Revert(migration); err != nil {
		return err
	}

	return nil
}

func newCmd(args arguments) error {
	if _, err := os.Stat(args.src); os.IsNotExist(err) {
		return err
	}

	if len(args.rest) < 2 {
		return errors.New("new requires a migration name given as an argument")
	}

	migration := gloat.GenerateMigration(strings.Join(args.rest[1:], "_"))
	migrationDirectoryPath := filepath.Join(args.src, migration.Path)

	if err := os.MkdirAll(migrationDirectoryPath, 0755); err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(migrationDirectoryPath, "up.sql"))
	if err != nil {
		return err
	}
	f.Close()

	f, err = os.Create(filepath.Join(migrationDirectoryPath, "down.sql"))
	if err != nil {
		return err
	}
	f.Close()

	fmt.Printf("Created %s\n", migrationDirectoryPath)

	return nil
}

func parseArguments() arguments {
	var args arguments

	urlDefault := os.Getenv("DATABASE_URL")
	urlUsage := `database connection url`

	srcDefault := os.Getenv("DATABASE_SRC")
	if srcDefault == "" {
		srcDefault = "db/migrations"
	}
	srcUsage := `the folder with migrations`

	flag.StringVar(&args.url, "url", urlDefault, urlUsage)
	flag.StringVar(&args.src, "src", srcDefault, srcUsage)

	flag.Usage = func() { fmt.Fprintf(os.Stderr, usage) }

	flag.Parse()

	args.rest = flag.Args()

	return args
}

func setupGloat(args arguments) (*gloat.Gloat, error) {
	u, err := url.Parse(args.url)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open(u.Scheme, args.url)
	if err != nil {
		return nil, err
	}

	store, err := databaseStoreFactory(u.Scheme, db)
	if err != nil {
		return nil, err
	}

	return &gloat.Gloat{
		Store:    store,
		Source:   gloat.NewFileSystemSource(args.src),
		Executor: gloat.NewSQLExecutor(db),
	}, nil
}

func databaseStoreFactory(driver string, db *sql.DB) (gloat.Store, error) {
	switch driver {
	case "postgres", "postgresql":
		return gloat.NewPostgreSQLStore(db), nil
	case "mysql":
		return gloat.NewMySQLStore(db), nil
	case "sqlite", "sqlite3":
		return gloat.NewMySQLStore(db), nil
	}

	return nil, errors.New("unsupported database driver " + driver)
}
