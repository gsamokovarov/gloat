package storage

import (
	"database/sql"

	"github.com/gsamokovarov/gloat/migration"
)

type DatabaseStorage struct {
	db *sql.DB

	createTableStatement         string
	insertMigrationStatement     string
	selectAllMigrationsStatement string
}

func (s *DatabaseStorage) Insert(migration *migration.Migration) error {
	if err := s.ensureSchemaTableExists(s.db); err != nil {
		return err
	}

	_, err := s.db.Exec(s.insertMigrationStatement, migration.Version)
	return err
}

func (s *DatabaseStorage) All() (migration.Migrations, error) {
	if err := s.ensureSchemaTableExists(s.db); err != nil {
		return nil, err
	}

	rows, err := s.db.Query(s.selectAllMigrationsStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var migrations migration.Migrations

	for rows.Next() {
		var m migration.Migration
		if err := rows.Scan(&m.Version); err != nil {
			return nil, err
		}

		migrations = append(migrations, m)
	}

	return migrations, nil
}

func (s *DatabaseStorage) ensureSchemaTableExists() error {
	_, err := db.Exec(s.createTableStatement)
	return err
}

func NewPostgresStorage(db *sql.DB) Storage {
	return &DatabaseStorage{
		db: db,
		createTableStatement: `CREATE TABLE IF NOT EXISTS schema_migrations (
			version SERIAL PRIMARY KEY)`,
		insertMigrationStatement:     `INSERT INTO schema_migrations (version) VALUES ($1)`,
		selectAllMigrationsStatement: `SELECT version FROM schema_migrations`,
	}
}

func NewMySQLStorage(db *sql.DB) Storage {
	return &DatabaseStorage{
		db: db,
		createTableStatement: `CREATE TABLE IF NOT EXISTS schema_migrations (
			version SERIAL PRIMARY KEY NOT NULL)`,
		insertMigrationStatement:     `INSERT INTO schema_migrations (version) VALUES ($1)`,
		selectAllMigrationsStatement: `SELECT version FROM schema_migrations`,
	}
}

func NewSQLite3Storage(db *sql.DB) Storage {
	return &DatabaseStorage{
		db: db,
		createTableStatement: `CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY AUTOINCREMENT)`,
		insertMigrationStatement:     `INSERT INTO schema_migrations (version) VALUES ($1)`,
		selectAllMigrationsStatement: `SELECT version FROM schema_migrations`,
	}
}
