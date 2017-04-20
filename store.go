package gloat

import (
	"database/sql"
	"errors"
)

type Store interface {
	Source

	Insert(*Migration) error
	Remove(*Migration) error
}

type DatabaseStore struct {
	db *sql.DB

	createTableStatement         string
	insertMigrationStatement     string
	removeMigrationStatement     string
	selectAllMigrationsStatement string
}

func (s *DatabaseStore) Insert(migration *Migration) error {
	if err := s.ensureSchemaTableExists(); err != nil {
		return err
	}

	_, err := s.db.Exec(s.insertMigrationStatement, migration.Version)
	return err
}

func (s *DatabaseStore) Remove(migration *Migration) error {
	if err := s.ensureSchemaTableExists(); err != nil {
		return err
	}

	_, err := s.db.Exec(s.removeMigrationStatement, migration.Version)
	return err
}

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

func NewPostgresSQLStore(db *sql.DB) Store {
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

func NewDatabaseStore(driver string, db *sql.DB) (Store, error) {
	switch driver {
	case "postgres":
		return NewPostgresSQLStore(db), nil
	case "mysql":
		return NewMySQLStore(db), nil
	case "sqlite", "sqlite3":
		return NewMySQLStore(db), nil
	}

	return nil, errors.New("unsupported database driver " + driver)
}
