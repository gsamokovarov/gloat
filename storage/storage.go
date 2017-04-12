package storage

import "github.com/gsamokovarov/gloat/migration"

type Storage interface {
	Insert(*migration.Migration) error
	Remove(*migration.Migration) error
	All() (migration.Migrations, error)
}
