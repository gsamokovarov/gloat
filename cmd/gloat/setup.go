package main

import (
	"database/sql"
	"errors"
	"os"

	"github.com/gsamokovarov/gloat"

	_ "github.com/lib/pq"
)

func setupGloat() (*gloat.Gloat, error) {
	connectionString, found := os.LookupEnv("DATABASE_URL")
	if !found {
		return nil, errors.New("no database config at DATABASE_URL")
	}

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	storage, err := gloat.NewDatabaseStorage("postgres", db)
	if err != nil {
		return nil, err
	}

	return &gloat.Gloat{
		InitialPath: "testdata/migrations",

		Storage:  storage,
		Source:   gloat.NewFileSystemSource("testdata/migrations"),
		Executor: gloat.NewExecutor(db),
	}, nil
}
