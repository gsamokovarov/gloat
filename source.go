package gloat

import (
	"os"
	"path/filepath"
)

type Source interface {
	Collect() (Migrations, error)
}

type FileSystemSource struct {
	MigrationsFolder string
}

func (s *FileSystemSource) Collect() (migrations Migrations, err error) {
	err = filepath.Walk(s.MigrationsFolder, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && path != s.MigrationsFolder {
			migration, err := FromPath(path)
			if err != nil {
				return err
			}

			migrations = append(migrations, migration)
		}

		return nil
	})

	return
}

func NewFileSystemSource(migrationsFolder string) Source {
	return &FileSystemSource{MigrationsFolder: migrationsFolder}
}
