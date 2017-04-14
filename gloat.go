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
	Executor *Executor
}

// UnappliedMigrations returns the unapplied migrations in the current
// gloat.
func (c *Gloat) UnappliedMigrations() (Migrations, error) {
	return UnappliedMigrations(c.Source, c.Storage)
}

// Apply applies a migrations.
func (c *Gloat) Apply(migration *Migration) error {
	return c.Executor.Up(migration, c.Storage)
}

// Reverse reverses a migrations.
func (c *Gloat) Revert(migration *Migration) error {
	return c.Executor.Down(migration, c.Storage)
}
