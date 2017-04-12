package source

import (
	"os"
	"path/filepath"

	"github.com/gsamokovarov/gloat/migration"
)

type FileSystemSource struct {
	MigrationsFolder string
}

func (s *FileSystemSource) Collect() (migration.Migrations, error) {
	var migrations migration.Migrations

	err := filepath.Walk(s.MigrationsFolder, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			migration, err := migration.FromPath(path)
			if err != nil {
				return err
			}

			migrations = append(migrations, migration)
		}

		return nil
	})

	return migrations, err
}
