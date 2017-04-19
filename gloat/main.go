package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"

	"github.com/gsamokovarov/gloat"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

const helpMsg = `Usage gloat: [COMMAND ...]

Gloat is a Go SQL migration utility.

Commands:
  up            Apply new migrations
  down          Revert the last applied migration

Options:
  -src          The folder with migrations
                (default $DATABASE_MIGRATIONS" or db/migrations)
  -url          The database connection URL
                (default $DATABASE_URL)
  -help         Show this message
`

var (
	urlFlag string
	srcFlag string
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, helpMsg)
		os.Exit(1)
	}

	args := parseArguments()

	var cmdName string
	if len(args) > 0 {
		cmdName = args[0]
	}

	var err error
	switch cmdName {
	case "up":
		err = upCmd()
	case "down":
		err = downCmd()
	default:
		fmt.Fprintf(os.Stderr, helpMsg)
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func upCmd() error {
	gl, err := setupGloat()
	if err != nil {
		return err
	}

	migrations, err := gl.Unapplied()
	if err != nil {
		return err
	}

	appliedMigrations := map[int64]bool{}

	for _, migration := range migrations {
		fmt.Printf("Applying migration: %d...\n", migration.Version)

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

func downCmd() error {
	gl, err := setupGloat()
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

	fmt.Printf("Reverting migration: %d...\n", migration.Version)

	if err := gl.Revert(migration); err != nil {
		return err
	}

	return nil
}

func parseArguments() []string {
	urlDefault := os.Getenv("DATABASE_URL")
	urlUsage := `database connection url`
	flag.StringVar(&urlFlag, "url", urlDefault, urlUsage)

	srcDefault := os.Getenv("DATABASE_MIGRATIONS")
	if srcDefault == "" {
		srcDefault = "db/migrations"
	}
	srcUsage := `the folder with migrations`
	flag.StringVar(&srcFlag, "src", srcDefault, srcUsage)

	flag.Parse()

	return flag.Args()
}

func setupGloat() (*gloat.Gloat, error) {
	if urlFlag == "" {
		return nil, errors.New("no -url or $DATABASE_URL present")
	}

	driver, err := guessDriver()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open(driver, urlFlag)
	if err != nil {
		return nil, err
	}

	storage, err := gloat.NewDatabaseStorage(driver, db)
	if err != nil {
		return nil, err
	}

	return &gloat.Gloat{
		InitialPath: srcFlag,

		Storage:  storage,
		Source:   gloat.NewFileSystemSource(srcFlag),
		Executor: gloat.NewExecutor(db),
	}, nil
}

func guessDriver() (string, error) {
	parsed, err := url.Parse(urlFlag)
	if err != nil {
		return "", err
	}

	return parsed.Scheme, nil
}
