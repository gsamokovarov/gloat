package gloat

import (
	"database/sql"
	"errors"
	"net/url"
	"os"
	"strings"
	"testing"

	// Needed to establish database connections during testing.
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var (
	gl Gloat

	db       *sql.DB
	dbURL    string
	dbSrc    string
	dbDriver string
)

type testingStore struct{ applied Migrations }

func (s *testingStore) Collect() (Migrations, error)                     { return s.applied, nil }
func (s *testingStore) Insert(migration *Migration, _ StoreExecer) error { return nil }
func (s *testingStore) Remove(migration *Migration, _ StoreExecer) error { return nil }

type testingExecutor struct{}

func (e *testingExecutor) Up(*Migration, Store) error   { return nil }
func (e *testingExecutor) Down(*Migration, Store) error { return nil }

type stubbedExecutor struct {
	up   func(*Migration, Store) error
	down func(*Migration, Store) error
}

func (e *stubbedExecutor) Up(m *Migration, s Store) error {
	if e.up != nil {
		return e.up(m, s)
	}

	return nil
}

func (e *stubbedExecutor) Down(m *Migration, s Store) error {
	if e.down != nil {
		e.down(m, s)
	}

	return nil
}

func cleanState(fn func()) error {
	_, err := db.Exec(`
		DROP TABLE IF EXISTS schema_migrations;	
		DROP TABLE IF EXISTS users;	
	`)

	if err != nil {
		return err
	}

	fn()

	return nil
}

func databaseStoreFactory(driver string, db *sql.DB) (Store, error) {
	switch driver {
	case "postgres", "postgresql":
		return NewPostgreSQLStore(db), nil
	case "mysql":
		return NewMySQLStore(db), nil
	case "sqlite", "sqlite3":
		return NewMySQLStore(db), nil
	}

	return nil, errors.New("unsupported database driver " + driver)
}

func TestUnapplied(t *testing.T) {
	gl.Store = &testingStore{applied: Migrations{}}

	migrations, err := gl.Unapplied()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if migrations[0].Version != 20170329154959 {
		t.Fatalf("Expected version 20170329154959, got: %d", migrations[0].Version)
	}
}

func TestUnapplied_Empty(t *testing.T) {
	gl.Store = &testingStore{
		applied: Migrations{
			&Migration{Version: 20170329154959},
		},
	}

	migrations, err := gl.Unapplied()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(migrations) != 0 {
		t.Fatalf("Expected no unapplied migrations, got: %v", migrations)
	}
}

func TestCurrent(t *testing.T) {
	gl.Store = &testingStore{
		applied: Migrations{
			&Migration{Version: 20170329154959},
		},
	}

	migration, err := gl.Current()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if migration == nil {
		t.Errorf("Expected current migration, got: %v", migration)
	}

	if migration.Version != 20170329154959 {
		t.Fatalf("Expected migration version to be 20170329154959, got: %d", migration.Version)
	}
}

func TestCurrent_Nil(t *testing.T) {
	gl.Store = &testingStore{}

	migration, err := gl.Current()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if migration != nil {
		t.Fatalf("Expected no current migration, got: %v", migration)
	}
}

func TestApply(t *testing.T) {
	called := false

	gl.Store = &testingStore{}
	gl.Executor = &stubbedExecutor{
		up: func(*Migration, Store) error {
			called = true
			return nil
		},
	}

	gl.Apply(nil)

	if !called {
		t.Fatalf("Expected Apply to call Executor.Up")
	}
}

func TestRevert(t *testing.T) {
	called := false

	gl.Store = &testingStore{}
	gl.Executor = &stubbedExecutor{
		down: func(*Migration, Store) error {
			called = true
			return nil
		},
	}

	gl.Revert(nil)

	if !called {
		t.Fatalf("Expected Revert to call Executor.Down")
	}
}

func init() {
	gl = Gloat{
		Source:   NewFileSystemSource("testdata/migrations"),
		Executor: &testingExecutor{},
	}

	dbURL = os.Getenv("DATABASE_URL")
	dbSrc = os.Getenv("DATABASE_SRC")

	{
		u, err := url.Parse(dbURL)
		if err != nil {
			panic(err)
		}

		dbDriver = u.Scheme
	}

	// Do a bit of post-processing so we can connect to non-postgres databases.
	if dbDriver != "postgres" {
		parts := strings.SplitN(dbURL, "://", 2)

		if len(parts) != 2 {
			panic("Cannot split " + dbURL + " into parts")
		}

		dbURL = parts[1]
	}

	{
		var err error

		db, err = sql.Open(dbDriver, dbURL)
		if err != nil {
			panic(err)
		}

		if err := db.Ping(); err != nil {
			panic(err)
		}
	}
}
