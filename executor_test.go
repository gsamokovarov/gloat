package gloat

import (
	"path/filepath"
	"testing"
)

func TestExecutorUp(t *testing.T) {
	td := filepath.Join(dbSrc, "20170329154959_introduce_domain_model")

	exe := NewExecutor(db)

	migration, err := MigrationFromPath(td)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	cleanState(func() {
		if err := exe.Up(migration, new(testingStore)); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		_, err := db.Exec(`SELECT id FROM users LIMIT 1`)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})
}

func TestExecutorDown(t *testing.T) {
	td := filepath.Join(dbSrc, "20170329154959_introduce_domain_model")

	exe := NewExecutor(db)

	migration, err := MigrationFromPath(td)
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

		_, err := db.Exec(`SELECT id FROM users LIMIT 1`)
		if err == nil {
			t.Error("Expected table users to be dropped")
		}
	})
}

func cleanState(fn func()) error {
	_, err := db.Exec(`
		DROP TABLE IF EXISTS schema_migrations;	
		DROP TABLE IF EXISTS users;	
	`)

	if err != nil {
		return err
	}

	fn()

	return nil
}
