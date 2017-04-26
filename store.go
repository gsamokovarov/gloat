package gloat

import (
	"database/sql"
)

// Store is an interface representing a place where the applied migrations are
// recorded.
type Store interface {
	Source

	Insert(*Migration) error
	Remove(*Migration) error
}

// DatabaseStore is a Store that keeps the applied migrations in a database
// table called schema_migrations. The table is automatically created if it
// does not exist.
type DatabaseStore struct {
	db *sql.DB

	createTableStatement         string
	insertMigrationStatement     string
	removeMigrationStatement     string
	selectAllMigrationsStatement string
}

// Insert records a migration version into the schema_migrations table.
func (s *DatabaseStore) Insert(migration *Migration) error {
	if err := s.ensureSchemaTableExists(); err != nil {
		return err
	}

	_, err := s.db.Exec(s.insertMigrationStatement, migration.Version)
	return err
}

// Remove removes a migration version from the schema_migrations table.
func (s *DatabaseStore) Remove(migration *Migration) error {
	if err := s.ensureSchemaTableExists(); err != nil {
		return err
	}

	_, err := s.db.Exec(s.removeMigrationStatement, migration.Version)
	return err
}

// Collect builds a slice of migrations with the versions of the recorded
// applied migrations.
func (s *DatabaseStore) Collect() (migrations Migrations, err error) {
	if err = s.ensureSchemaTableExists(); err != nil {
		return
	}

	rows, err := s.db.Query(s.selectAllMigrationsStatement)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		migration := &Migration{}
		if err = rows.Scan(&migration.Version); err != nil {
			return
		}

		migrations = append(migrations, migration)
	}

	return
}

func (s *DatabaseStore) ensureSchemaTableExists() error {
	_, err := s.db.Exec(s.createTableStatement)
	return err
}

// NewPostgreSQLStore creates a Store for PostgreSQL.
func NewPostgreSQLStore(db *sql.DB) Store {
	return &DatabaseStore{
		db: db,
		createTableStatement: `
			CREATE TABLE IF NOT EXISTS schema_migrations (
				version BIGINT PRIMARY KEY NOT NULL
			)`,
		insertMigrationStatement: `
			INSERT INTO schema_migrations (version)
			VALUES ($1)`,
		removeMigrationStatement: `
			DELETE FROM schema_migrations
			WHERE version=$1`,
		selectAllMigrationsStatement: `
			SELECT version
			FROM schema_migrations`,
	}
}

// NewMySQLStore creates a Store for MySQL.
func NewMySQLStore(db *sql.DB) Store {
	return &DatabaseStore{
		db: db,
		createTableStatement: `
			CREATE TABLE IF NOT EXISTS schema_migrations (
				version BIGINT PRIMARY KEY NOT NULL
			)`,
		insertMigrationStatement: `
			INSERT INTO schema_migrations (version)
			VALUES (?)`,
		removeMigrationStatement: `
			DELETE FROM schema_migrations
			WHERE version=?`,
		selectAllMigrationsStatement: `
			SELECT version
			FROM schema_migrations`,
	}
}

// NewSQLite3Store creates a Store for SQLite3.
func NewSQLite3Store(db *sql.DB) Store {
	return &DatabaseStore{
		db: db,
		createTableStatement: `
			CREATE TABLE IF NOT EXISTS schema_migrations (
				version BIGINT PRIMARY KEY NOT NULL
			)`,
		insertMigrationStatement: `
			INSERT INTO schema_migrations (version)
			VALUES (?)`,
		removeMigrationStatement: `
			DELETE FROM schema_migrations
			WHERE version=?`,
		selectAllMigrationsStatement: `
			SELECT version
			FROM schema_migrations`,
	}
}
