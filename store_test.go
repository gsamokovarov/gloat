package gloat

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestDatabaseStore_Insert(t *testing.T) {
	td := filepath.Join(dbSrc, "20170329154959_introduce_domain_model")

	migration, err := MigrationFromPath(td)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	dbStore, err := NewDatabaseStore(dbDriver, db)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	cleanState(func() {
		if _, err := db.Exec(`SELECT version FROM schema_migrations`); err == nil {
			t.Fatal("Expected table schema_migrations to not exist")
		}

		if err := dbStore.Insert(migration); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		var version int64

		err := db.QueryRow(`SELECT version FROM schema_migrations`).Scan(&version)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if version != 20170329154959 {
			t.Fatalf("Expected version to be 20170329154959, got %d", version)
		}
	})
}

func TestDatabaseStore_Remove(t *testing.T) {
	td := filepath.Join(dbSrc, "20170329154959_introduce_domain_model")

	migration, err := MigrationFromPath(td)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	dbStore, err := NewDatabaseStore(dbDriver, db)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	cleanState(func() {
		if err := dbStore.Insert(migration); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if err := dbStore.Remove(migration); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		var version int64

		row := db.QueryRow(`SELECT version FROM schema_migrations`)
		if row.Scan(&version) == nil {
			t.Fatal("Expected error no rows in result set")
		}
	})
}

func TestDatabaseStore_Collect(t *testing.T) {
	td := filepath.Join(dbSrc, "20170329154959_introduce_domain_model")

	migration, err := MigrationFromPath(td)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	dbStore, err := NewDatabaseStore(dbDriver, db)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	cleanState(func() {
		if err := dbStore.Insert(migration); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		migrations, err := dbStore.Collect()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		expectedMigrations := Migrations{
			&Migration{Version: 20170329154959},
		}

		if !reflect.DeepEqual(migrations, expectedMigrations) {
			t.Fatalf("Expected migrations to be: %v, got %v", expectedMigrations, migrations)
		}
	})
}
