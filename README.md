<p align=right>
	<a href="https://travis-ci.org/gsamokovarov/gloat">
		<img src="https://travis-ci.org/gsamokovarov/gloat.svg?branch=master" alt="Build Status" data-canonical-src="https://travis-ci.org/gsamokovarov/gloat.svg?branch=master">
	</a>
</p>

# Gloat /ɡlōt/

> Contemplate or dwell on one's own success or another's misfortune with
> smugness or malignant pleasure.

Gloat is a modular SQL migration library for the Go programming language. Being
a library, gloat can be easily integrated into your application or ORM.

## Library

If you are using gloat as a library, the main components you'll be dealing with
are migration, source, store and SQL executor. You'll be using those through the
methods on the `Gloat` struct.

```go
db, err := sql.Open("postgres", "connection string")
if err != nil {
	// Handle the *sql.DB creation error.
}

gl := gloat.Gloat{
	Store:    gloat.NewPostgreSQLStore(db),
	Source:   gloat.NewFileSystemSource("migrations"),
	Executor: gloat.NewSQLExecutor(db),
}
```

### Migration

Migration holds all the relevant information for a migration. The content of
the forward (up) side of a migration, the backward (down) side, a path and
version. The version is used to determine the order of which the migrations
would be executed. The path is the name in a store.

```go
type Migration struct {
	UpSQL   []byte
	DownSQL []byte
	Path    string
	Version int64
}
```

### Source

The `Source` interface represents a source of migration. The most common source
is the file system.

```go
type Source interface {
	Collect() (Migrations, error)
}
```


`gloat.NewFileSystemSource` is a constructor function that creates a source
that collects migrations from a folder with the following structure:

```
migrations/
└── 20170329154959_introduce_domain_model
    ├── down.sql
    └── up.sql
```

In the example above `migrations` is a folder that stores all of the
migrations. A migrations itself is a folder with a name in the form of
`:timestamp_:name` containing `up.sql` and `down.sql` files.

The `up.sql` file contains the SQL that's executed when a migration is applied:

```sql
CREATE TABLE users (
    id    bigserial PRIMARY KEY NOT NULL,
    name  character varying NOT NULL,
    email character varying NOT NULL,

    created_at  timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at  timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL
);
```

While `down.sql` is executed when a migration is reverted:

```sql
DROP TABLE users;
```

If the `down.sql` file is not present, we say that a migration is irreversible.

## Store

The Store is an interface representing a place where the applied migrations are
recorded. The common thing is to store the migrations in a database table. The
`gloat.DatabaseStore` does just that. The `gloat.NewPostgreSQLStore` constructor
function creates such store that records the migration in a table called
`schema_migrations`. The table is automatically created if it does not exist.

```go
type Store interface {
	Source

	Insert(*Migration) error
	Remove(*Migration) error
}
```

The `Store.Insert` records the migration version in to the `schema_migrations`
table, while `Store.Remove` deletes the column with the version from
the table. There are the following builtin store constructors:

```go
// NewPostgreSQLStore creates a Store for PostgreSQL.
func NewPostgreSQLStore(db *sql.DB) Store {}

// NewMySQLStore creates a Store for MySQL.
func NewMySQLStore(db *sql.DB) Store {}

// NewSQLite3Store creates a Store for SQLite3.
func NewSQLite3Store(db *sql.DB) Store {}
```

### Executor

The `Executor` interface, well, it executes the migrations. For SQL migrations,
there is the `gloat.SQLExecutor` implementation. It's an interface,
nevertheless, so you can fake it out during testing.

```go
type Executor interface {
	Up(*Migration, Store) error
	Down(*Migration, Store) error
}
```

The executor executes the migration `UpSQL` or `DownSQL` sections.

### Gloat

A `Gloat` binds a migration `Source`, `Store` and `Executor` into one thing, so
it's easier to `Apply`, and `Revert` migrations.

```go
gl := gloat.Gloat{
	Store:    gloat.NewPostgreSQLStore(db),
	Source:   gloat.NewFileSystemSource("migrations"),
	Executor: gloat.NewSQLExecutor(db),
}

// Applies all of the unapplied migrations.
if migrations, err := gl.Unapplied(); err == nil {
	for _, migration := range migrations {
		gl.Apply(migration)
	}
}

// Revert the last applied migration.
if migration, err := gl.Current(); err == nil {
	gl.Revert(migration)
}
```

Here is a description for the main Gloat methods.

```go
// Unapplied returns the unapplied migrations in the current gloat.
func (c *Gloat) Unapplied() (Migrations, error) {}

// Current returns the latest applied migration. Even if no error is returned,
// the current migration can be nil.
//
// This is the case when the last applied migration is no longer available from
// the source or there are no migrations to begin with.
func (c *Gloat) Current() (*Migration, error) {}

// Apply applies a migration.
func (c *Gloat) Apply(migration *Migration) error {}

// Revert rollbacks a migration.
func (c *Gloat) Revert(migration *Migration) error {}
```
