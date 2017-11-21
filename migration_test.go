package gloat

import (
	"io/ioutil"
	"testing"

	"github.com/gsamokovarov/assert"
)

func TestMigrationReversible(t *testing.T) {
	m := Migration{}

	assert.False(t, m.Reversible())

	m.DownSQL = []byte("DROP TABLE users;")

	assert.True(t, m.Reversible())
}

func TestMigrationPersistable(t *testing.T) {
	m := Migration{}

	if m.Persistable() {
		t.Fatalf("Expected %v to not be persistable", m)
	}

	m.Path = "migrations/0001_something"
	assert.True(t, m.Persistable())
}

func TestMigrationFromPath(t *testing.T) {
	expectedPath := "testdata/migrations/20170329154959_introduce_domain_model"

	m, err := MigrationFromBytes(expectedPath, ioutil.ReadFile)
	assert.Nil(t, err)

	assert.Equal(t, 20170329154959, m.Version)
	assert.Equal(t, expectedPath, m.Path)
}

func TestMigrationsExcept(t *testing.T) {
	var migrations Migrations

	expectedPath := "testdata/migrations/20170329154959_introduce_domain_model"

	m, err := MigrationFromBytes(expectedPath, ioutil.ReadFile)
	assert.Nil(t, err)

	migrations = append(migrations, m)

	exceptedMigrations := migrations.Except(nil)
	assert.Equal(t, exceptedMigrations, migrations)

	exceptedMigrations = migrations.Except(Migrations{m})
	assert.Len(t, 0, exceptedMigrations)
}
