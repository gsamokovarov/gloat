package source

import "github.com/gsamokovarov/gloat/migration"

type Source interface {
	Collect() (migration.Migrations, error)
}
