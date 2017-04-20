package gloat

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestFileSystemSourceCollect(t *testing.T) {
	td := "testdata/migrations"
	fs := NewFileSystemSource(td)

	migrations, err := fs.Collect()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	m, err := MigrationFromPath(filepath.Join(td, "20170329154959_introduce_domain_model"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expectedMigrations := Migrations{m}

	if !reflect.DeepEqual(migrations, expectedMigrations) {
		t.Fatalf("Expected migrations to be: %v, got %v", expectedMigrations, migrations)
	}
}

func TestFileSystemSourceCollectEmpty(t *testing.T) {
	td := "testdata/no_migrations"
	fs := NewFileSystemSource(td)

	migrations, err := fs.Collect()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(migrations) != 0 {
		t.Fatalf("Expected no migrations collected in: %s", td)
	}
}
