package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/gsamokovarov/gloat"
	_ "github.com/lib/pq"
)

func main() {
	connectionString, found := os.LookupEnv("DATABASE_URL")
	if !found {
		Exitf(1, "No database config at DATABASE_URL")
	}

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		Exitf(1, "Error: %v\n", err)
	}

	executor := gloat.NewExecutor(db)

	source := gloat.NewFileSystemSource("testdata/migrations")
	storage := gloat.NewPostgresSQLStorage(db)

	migrations, err := gloat.UnappliedMigrations(source, storage)
	if err != nil {
		Exitf(1, "Error: %v\n", err)
	}

	appliedMigrations := map[int64]bool{}

	for _, migration := range migrations {
		Outf("Applying migration: %d...\n", migration.Version)

		if err := executor.Up(migration, storage); err != nil {
			Exitf(1, "Error: %v\n", err)
		}

		appliedMigrations[migration.Version] = true
	}

	if len(appliedMigrations) == 0 {
		Outf("No migrations to apply\n")
	}
}

func Exitf(code int, format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(code)
}

func Outf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stdout, format, a...)
}
