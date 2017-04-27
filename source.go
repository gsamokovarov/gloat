package gloat

import (
	"io/ioutil"
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
	Dir string
}

// Collect builds migrations stored in a folder like the following structure:
//
// migrations/
// └── 20170329154959_introduce_domain_model
//     ├── down.sql
//     └── up.sql
func (s *FileSystemSource) Collect() (migrations Migrations, err error) {
	err = filepath.Walk(s.Dir, func(path string, info os.FileInfo, err error) error {
		if info != nil && info.IsDir() && path != s.Dir {
			migration, err := MigrationFromBytes(path, ioutil.ReadFile)
			if err != nil {
				return err
			}

			migrations = append(migrations, migration)
		}

		return nil
	})

	migrations.Sort()

	return
}

// NewFileSystemSource creates a new source of migrations that takes them right
// out of the file system.
func NewFileSystemSource(dir string) Source {
	return &FileSystemSource{Dir: dir}
}

// AssetSource is a go-bindata migration source for binary embedded migrations.
// You need to pass a prefix, the Asset and AssetDir functions, go-bindata
// generates.
type AssetSource struct {
	Prefix   string
	Asset    func(string) ([]byte, error)
	AssetDir func(string) ([]string, error)
}

// Collect builds migrations from a go-bindata embedded migrations.
func (s *AssetSource) Collect() (migrations Migrations, err error) {
	dirs, err := s.AssetDir(s.Prefix)
	if err != nil {
		return
	}

	for _, path := range dirs {
		var migration *Migration

		migration, err = MigrationFromBytes(filepath.Join(s.Prefix, path), s.Asset)
		if err != nil {
			return
		}

		migrations = append(migrations, migration)
	}

	migrations.Sort()

	return
}

// NewAssetSource creates a new source of binary migrations with go-bindata.
func NewAssetSource(prefix string, asset func(string) ([]byte, error), assetDir func(string) ([]string, error)) Source {
	return &AssetSource{Prefix: prefix, Asset: asset, AssetDir: assetDir}
}
