package gloat

import (
	"database/sql"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/gsamokovarov/assert"
)

func TestDatabaseStore_Insert(t *testing.T) {
	td := filepath.Join(dbSrc, "20170329154959_introduce_domain_model")

	migration, err := MigrationFromBytes(td, ioutil.ReadFile)
	assert.Nil(t, err)

	dbStore, err := databaseStoreFactory(dbDriver, db)
	assert.Nil(t, err)

	cleanState(func() {
		_, err := db.Exec(`SELECT version FROM schema_migrations`)
		assert.NotNil(t, err)

		err = dbStore.Insert(migration, nil)
		assert.Nil(t, err)

		var version int64

		err = db.QueryRow(`SELECT version FROM schema_migrations`).Scan(&version)
		assert.Nil(t, err)

		assert.Equal(t, 20170329154959, version)
	})
}

func TestDatabaseStore_Remove(t *testing.T) {
	td := filepath.Join(dbSrc, "20170329154959_introduce_domain_model")

	migration, err := MigrationFromBytes(td, ioutil.ReadFile)
	assert.Nil(t, err)

	dbStore, err := databaseStoreFactory(dbDriver, db)
	assert.Nil(t, err)

	cleanState(func() {
		err := dbStore.Insert(migration, nil)
		assert.Nil(t, err)

		err = dbStore.Remove(migration, nil)
		assert.Nil(t, err)

		var version int64

		err = db.QueryRow(`SELECT version FROM schema_migrations`).Scan(&version)
		assert.Equal(t, sql.ErrNoRows, err)
	})
}

func TestDatabaseStore_Collect(t *testing.T) {
	td := filepath.Join(dbSrc, "20170329154959_introduce_domain_model")

	migration, err := MigrationFromBytes(td, ioutil.ReadFile)
	assert.Nil(t, err)

	dbStore, err := databaseStoreFactory(dbDriver, db)
	assert.Nil(t, err)

	cleanState(func() {
		err := dbStore.Insert(migration, nil)
		assert.Nil(t, err)

		migrations, err := dbStore.Collect()
		assert.Nil(t, err)

		expectedMigrations := Migrations{
			&Migration{Version: 20170329154959},
		}

		assert.Equal(t, migrations, expectedMigrations)
	})
}
