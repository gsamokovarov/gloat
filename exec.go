package gloat

import (
	"database/sql"
	"fmt"
)

// IrreversibleError is the error return when we're trying to reverse a
// migration that has a blank down SQL content.
type IrreversibleError struct {
	Version int64
}

// Error implements the error interface.
func (err IrreversibleError) Error() string {
	return fmt.Sprintf("cannot reverse migration %s", err.Version)
}

// Executor is a type that executes migrations up and down.
type Executor struct {
	db *sql.DB
}

// Up executes a migrations forward.
func (e *Executor) Up(migration *Migration, storage Storage) error {
	if _, err := e.db.Exec(string(migration.UpSQL)); err != nil {
		return err
	}

	if err := storage.Insert(migration); err != nil {
		return err
	}

	return nil
}

// Down reverses a migrations.
func (e *Executor) Down(migration *Migration, storage Storage) error {
	if !migration.Reversible() {
		return IrreversibleError{migration.Version}
	}

	if _, err := e.db.Exec(string(migration.DownSQL)); err != nil {
		return err
	}

	if err := storage.Remove(migration); err != nil {
		return err
	}

	return nil
}
