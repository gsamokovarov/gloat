package main

import (
	"database/sql"
	"os"

	"github.com/gsamokovarov/gloat"

	_ "github.com/lib/pq"
)

var gl gloat.Gloat

func init() {
	connectionString, found := os.LookupEnv("DATABASE_URL")
	if !found {
		Exitf(1, "No database config at DATABASE_URL")
	}

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		Exitf(1, "Error: %v\n", err)
	}

	gl = gloat.Gloat{
		InitialPath: "testdata/migrations",

		Source:   gloat.NewFileSystemSource("testdata/migrations"),
		Storage:  gloat.NewGenericDatabaseStorage(db),
		Executor: gloat.NewExecutor(db),
	}
}
