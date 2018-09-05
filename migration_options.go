package gloat

import (
	"bytes"
	"encoding/json"
)

// MigrationOptions are the options for a migration. Keep in mind that some
// options (transaction) are not supported by every RDBMS (ahem, MySQL).
type MigrationOptions struct {
	Transaction bool `json:"transaction"`
}

// DefaultMigrationOptions generate the default migration options.
//
// By default, the migrations are run in transaction, you can optionally
// configure them not to. For now, that is. 😅
func DefaultMigrationOptions() MigrationOptions {
	return MigrationOptions{
		Transaction: true,
	}
}

func parseMigrationOptions(data []byte) (options MigrationOptions, err error) {
	if data == nil {
		return DefaultMigrationOptions(), nil
	}

	err = json.NewDecoder(bytes.NewReader(data)).Decode(&options)

	return
}
