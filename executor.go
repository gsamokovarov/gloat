package gloat

import (
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
	Up(*Migration, Store) error
	Down(*Migration, Store) error
}

// SQLExecutor is a type that executes migrations in a database.
type SQLExecutor struct {
	db SQLTransactor
}

// Up applies a migration.
func (e *SQLExecutor) Up(migration *Migration, store Store) error {
	return e.exec(migration.Options.Transaction, func(tx SQLExecer) error {
		if _, err := tx.Exec(string(migration.UpSQL)); err != nil {
			return err
		}

		return store.Insert(migration, tx)
	})
}

// Down reverses a migrations.
func (e *SQLExecutor) Down(migration *Migration, store Store) error {
	if !migration.Reversible() {
		return IrreversibleError{migration.Version}
	}

	return e.exec(migration.Options.Transaction, func(tx SQLExecer) error {
		if _, err := tx.Exec(string(migration.DownSQL)); err != nil {
			return err
		}

		return store.Remove(migration, tx)
	})
}

func (e *SQLExecutor) exec(transaction bool, action func(SQLExecer) error) error {
	if !transaction {
		return action(e.db)
	}

	tx, err := e.db.Begin()
	if err != nil {
		return err
	}

	if err := action(tx); err != nil {
		defer tx.Rollback()
		return err
	}

	return tx.Commit()
}

// NewSQLExecutor creates an SQLExecutor.
func NewSQLExecutor(db SQLTransactor) Executor {
	return &SQLExecutor{db: db}
}
