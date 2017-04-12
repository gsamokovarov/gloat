package exec

import (
	"database/sql"
	"fmt"

	"github.com/gsamokovarov/gloat/migration"
	"github.com/gsamokovarov/gloat/source"
	"github.com/gsamokovarov/gloat/storage"
)

// UnappliedMigrations selects the unapplied migrations from a Source. For a
// migration to be unapplied it should not be present in the Storage.
func UnappliedMigrations(source source.Source, storage source.Storage) (migration.Migrations, error) {
	allMigrations, err := source.Collect()
	if err != nil {
		return nil, err
	}

	appliedMigrations, err := storage.All()
	if err != nil {
		return nil, err
	}

	return allMigrations.Filter(appliedMigrations), nil
}

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
func (e *Executor) Up(migration migration.Migration, storage storage.Storage) error {
	if _, err := db.Exec(migration.upSQL); err != nil {
		return err
	}

	if err := storage.Insert(migration); err != nil {
		return err
	}

	return nil
}

// Down reverses a migrations.
func (e *Executor) Down(migration migration.Migration, storage storage.Storage) error {
	if !migation.Reversible() {
		return IrreversibleError{migration.Version}
	}

	if _, err := db.Exec(migration.downSQL); err != nil {
		return err
	}

	if err := storage.Remove(migration); err != nil {
		return err
	}

	return nil
}
