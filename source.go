package gloat

import (
	"os"
	"path/filepath"
)

// Source is an interface representing a migrations source.
type Source interface {
	Collect() (Migrations, error)
}

// FileSystemSource is a file system source of migrations. The migrations are
// stored in folders with the following structure:
//
// migrations/
// └── 20170329154959_introduce_domain_model
//     ├── down.sql
//     └── up.sql
type FileSystemSource struct {
	MigrationsFolder string
}

// Collect builds migrations stored in a migrations folder like the following
// the one below:
//
// migrations/
// └── 20170329154959_introduce_domain_model
//     ├── down.sql
//     └── up.sql
func (s *FileSystemSource) Collect() (migrations Migrations, err error) {
	err = filepath.Walk(s.MigrationsFolder, func(path string, info os.FileInfo, err error) error {
		if info != nil && info.IsDir() && path != s.MigrationsFolder {
			migration, err := MigrationFromPath(path)
			if err != nil {
				return err
			}

			migrations = append(migrations, migration)
		}

		return nil
	})

	return
}

// NewFileSystemSource creates a new source of migrations that takes them right
// out of the file system.
func NewFileSystemSource(migrationsFolder string) Source {
	return &FileSystemSource{MigrationsFolder: migrationsFolder}
}
