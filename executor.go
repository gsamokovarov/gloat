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
	return fmt.Sprintf("cannot reverse migration %d", err.Version)
}

// Executor is a type that executes migrations up and down.
type Executor interface {
	Up(*Migration, Storage) error
	Down(*Migration, Storage) error
}

// Executor is a type that executes migrations in a database.
type SQLExecutor struct {
	db *sql.DB
}

// Up applies a migrations.
func (e *SQLExecutor) Up(migration *Migration, storage Storage) error {
	tx, err := e.db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(string(migration.UpSQL)); err != nil {
		tx.Rollback()
		return err
	}

	if err := storage.Insert(migration); err != nil {
		return err
	}

	return tx.Commit()
}

// Down reverses a migrations.
func (e *SQLExecutor) Down(migration *Migration, storage Storage) error {
	if !migration.Reversible() {
		return IrreversibleError{migration.Version}
	}

	tx, err := e.db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(string(migration.DownSQL)); err != nil {
		tx.Rollback()
		return err
	}

	if err := storage.Remove(migration); err != nil {
		return err
	}

	return tx.Commit()
}

func NewExecutor(db *sql.DB) Executor {
	return &SQLExecutor{db: db}
}
