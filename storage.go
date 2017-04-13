package gloat

import "database/sql"

type Storage interface {
	Source

	Insert(*Migration) error
	Remove(*Migration) error
}

type DatabaseStorage struct {
	db *sql.DB

	createTableStatement         string
	insertMigrationStatement     string
	removeMigrationStatement     string
	selectAllMigrationsStatement string
}

func (s *DatabaseStorage) Insert(migration *Migration) error {
	if err := s.ensureSchemaTableExists(); err != nil {
		return err
	}

	_, err := s.db.Exec(s.insertMigrationStatement, migration.Version)
	return err
}

func (s *DatabaseStorage) Remove(migration *Migration) error {
	if err := s.ensureSchemaTableExists(); err != nil {
		return err
	}

	_, err := s.db.Exec(s.removeMigrationStatement, migration.Version)
	return err
}

func (s *DatabaseStorage) Collect() (Migrations, error) {
	if err := s.ensureSchemaTableExists(); err != nil {
		return nil, err
	}

	rows, err := s.db.Query(s.selectAllMigrationsStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var migrations Migrations

	for rows.Next() {
		m := &Migration{}
		if err := rows.Scan(&m.Version); err != nil {
			return nil, err
		}

		migrations = append(migrations, m)
	}

	return migrations, nil
}

func (s *DatabaseStorage) ensureSchemaTableExists() error {
	_, err := s.db.Exec(s.createTableStatement)
	return err
}

func NewPostgresSQLStorage(db *sql.DB) Storage {
	return &DatabaseStorage{
		db: db,
		createTableStatement: `
			CREATE TABLE IF NOT EXISTS schema_migrations (
				version BIGSERIAL PRIMARY KEY
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

func NewMySQLStorage(db *sql.DB) Storage {
	return &DatabaseStorage{
		db: db,
		createTableStatement: `
			CREATE TABLE IF NOT EXISTS schema_migrations (
				version BIGINT PRIMARY KEY AUTO_INCREMENT NOT NULL
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

func NewSQLite3Storage(db *sql.DB) Storage {
	return &DatabaseStorage{
		db: db,
		createTableStatement: `
			CREATE TABLE IF NOT EXISTS schema_migrations (
				version BIGINT PRIMARY KEY AUTOINCREMENT NOT NULL
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
