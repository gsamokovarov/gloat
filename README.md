<p align=right>
	[![Build Status](https://travis-ci.org/gsamokovarov/gloat.svg?branch=master)](https://travis-ci.org/gsamokovarov/gloat)
</p>

# Gloat /ɡlōt/

> Contemplate or dwell on one's own success or another's misfortune with
> smugness or malignant pleasure.

Gloat is a modular SQL migration library for the Go programming language. Being
a library, gloat can be easily integrated into your application or ORM.

## Library

If you are using gloat as a library, there are 3 main components you need to
understand. A migration source, migration storage, migration executor and a
Gloat struct that binds them all together.

```go
db, err := sql.Open("postgres", "connection string")
if err != nil {
	// Handle the *sql.DB creation error.
}

store, _ := gloat.NewDatabaseStore("postgres", db)
if err != nil {
	// The supported RDBMS drivers are `postgresql`, `mysql` and `sqlite`. The
	// error indicates unsuported one.
}

gl := Gloat{
	Store:    store,
	Source:   gloat.NewFileSystemSource("migrations"),
	Executor: gloat.NewSQLExecutor(db),
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


`gloat.NewFileSystemSource` is a constructor function that creates a storage in
which collects migrations from a folder with the following structure:

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
`gloat.DatabaseStore` does just that. The `gloat.NewDatabaseStore` constructor
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
table, while `Store.Remove` deletes the column with the version from the table.

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
store, _ := gloat.NewDatabaseStore("postgres", db)
if err != nil {
	// The supported RDBMS drivers are `postgresql`, `mysql` and `sqlite`. The
	// error indicates unsuported one.
}

gl := Gloat{
	Store:    store,
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
