package gloat

import (
	"io/ioutil"
	"testing"

	"github.com/gsamokovarov/assert"
)

func TestDefaultMigrationOptions(t *testing.T) {
	options := DefaultMigrationOptions()

	assert.True(t, options.Transaction)
}

func TestMigrationWithoutExplicitOptions(t *testing.T) {
	expectedPath := "testdata/migrations/20170329154959_introduce_domain_model"

	m, err := MigrationFromBytes(expectedPath, ioutil.ReadFile)
	assert.Nil(t, err)

	assert.True(t, m.Options.Transaction)
}

func TestMigrationWithExplicitOptions(t *testing.T) {
	expectedPath := "testdata/migrations/20180905150724_concurrent_migration"

	m, err := MigrationFromBytes(expectedPath, ioutil.ReadFile)
	assert.Nil(t, err)

	assert.False(t, m.Options.Transaction)
}
