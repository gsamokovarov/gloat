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
		Exitf(1, "Error: %v", err)
	}

	executor := gloat.NewExecutor(db)

	source := gloat.NewFileSystemSource("testdata/migrations")
	storage := gloat.NewPostgresSQLStorage(db)

	migrations, err := gloat.UnappliedMigrations(source, storage)
	if err != nil {
		Exitf(1, "Error: %v", err)
	}

	for _, migration := range migrations {
		if err := executor.Up(migration, storage); err != nil {
			Exitf(1, "Error: %v", err)
		}
	}
}

func Exitf(code int, format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(code)
}
