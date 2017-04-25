package gloat

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestSQLExecutor_Up(t *testing.T) {
	td := filepath.Join(dbSrc, "20170329154959_introduce_domain_model")

	exe := NewSQLExecutor(db)

	migration, err := MigrationFromBytes(td, ioutil.ReadFile)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	cleanState(func() {
		if err := exe.Up(migration, new(testingStore)); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if _, err := db.Exec(`SELECT id FROM users LIMIT 1`); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})
}

func TestSQLExecutor_Down(t *testing.T) {
	td := filepath.Join(dbSrc, "20170329154959_introduce_domain_model")

	exe := NewSQLExecutor(db)

	migration, err := MigrationFromBytes(td, ioutil.ReadFile)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	cleanState(func() {
		if err := exe.Up(migration, new(testingStore)); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if err := exe.Down(migration, new(testingStore)); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if _, err := db.Exec(`SELECT id FROM users LIMIT 1`); err == nil {
			t.Error("Expected table users to be dropped")
		}
	})
}
