package gloat

// Gloat glues all the components needed to apply and revert
// migrations.
type Gloat struct {
	// InitialPath is the base path to the source of the migrations. If the
	// source is a File System, this can be the folder storing all of the
	// migrations.
	InitialPath string

	// Source is an incoming source of migrations. It can be File System or
	// embedded migrations with go-bindata, etc.
	Source Source

	// Storage is the place where we store the applied migration versions. Can
	// be one of the builtin database storages, etc.
	Storage Storage

	// Executor applies migrations and marks the newly applied migration
	// versions in the Store.
	Executor Executor
}

// Unapplied returns the unapplied migrations in the current gloat.
func (c *Gloat) Unapplied() (Migrations, error) {
	return UnappliedMigrations(c.Source, c.Storage)
}

// Current returns the latest applied migration. Even if no error is returned,
// the current migration can be nil.
//
// This is the case when the last applied migration is no longer available from
// the source or there are no migrations to begin with.
func (c *Gloat) Current() (*Migration, error) {
	appliedMigrations, err := c.Storage.Collect()
	if err != nil {
		return nil, err
	}

	currentMigration := appliedMigrations.Current()
	if currentMigration == nil {
		return nil, nil
	}

	availableMigrations, err := c.Source.Collect()
	if err != nil {
		return nil, err
	}

	for i := len(availableMigrations) - 1; i >= 0; i-- {
		migration := availableMigrations[i]

		if migration.Version == currentMigration.Version {
			return migration, nil
		}
	}

	return nil, nil
}

// Apply applies a migration.
func (c *Gloat) Apply(migration *Migration) error {
	return c.Executor.Up(migration, c.Storage)
}

// Revert rollbacks a migration.
func (c *Gloat) Revert(migration *Migration) error {
	return c.Executor.Down(migration, c.Storage)
}
