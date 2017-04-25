package gloat

import (
	"io/ioutil"
	"reflect"
	"testing"
)

func TestMigrationReversible(t *testing.T) {
	m := Migration{}

	if m.Reversible() {
		t.Fatalf("Expected %v to not be reversible", m)
	}

	m.DownSQL = []byte("DROP TABLE users;")

	if !m.Reversible() {
		t.Fatalf("Expected %v to be reversible", m)
	}
}

func TestMigrationPersistable(t *testing.T) {
	m := Migration{}

	if m.Persistable() {
		t.Fatalf("Expected %v to not be persistable", m)
	}

	m.Path = "migrations/0001_something"
	if !m.Persistable() {
		t.Fatalf("Expected %v to be persistable", m)
	}
}

func TestMigrationFromPath(t *testing.T) {
	expectedPath := "testdata/migrations/20170329154959_introduce_domain_model"

	m, err := MigrationFromBytes(expectedPath, ioutil.ReadFile)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if m.Version != 20170329154959 {
		t.Fatalf("Expected migration with version 20170329154959, got %d", m.Version)
	}

	if m.Path != expectedPath {
		t.Fatalf("Expected migration with path %s, got: %s", expectedPath, m.Path)
	}
}

func TestMigrationsExcept(t *testing.T) {
	var migrations Migrations

	expectedPath := "testdata/migrations/20170329154959_introduce_domain_model"

	m, err := MigrationFromBytes(expectedPath, ioutil.ReadFile)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	migrations = append(migrations, m)

	exceptedMigrations := migrations.Except(nil)
	if !reflect.DeepEqual(migrations, exceptedMigrations) {
		t.Fatalf("Expected exceptedMigrations to be unchanged, got: %v", exceptedMigrations)
	}

	exceptedMigrations = migrations.Except(Migrations{m})
	if len(exceptedMigrations) != 0 {
		t.Fatalf("Expected exceptedMigrations to be empty, got: %v", exceptedMigrations)
	}
}
