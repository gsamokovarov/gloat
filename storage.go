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

func (s *DatabaseStorage) Collect() (migrations Migrations, err error) {
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

func (s *DatabaseStorage) ensureSchemaTableExists() error {
	_, err := s.db.Exec(s.createTableStatement)
	return err
}

func NewGenericDatabaseStorage(db *sql.DB) Storage {
	return &DatabaseStorage{
		db: db,
		createTableStatement: `
			CREATE TABLE IF NOT EXISTS schema_migrations (
				version bigint PRIMARY KEY NOT NULL
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
