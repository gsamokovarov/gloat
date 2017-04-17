package gloat

import "testing"

var gl Gloat

type testingStorage struct{ applied Migrations }

func (s *testingStorage) Collect() (Migrations, error)      { return s.applied, nil }
func (s *testingStorage) Insert(migration *Migration) error { return nil }
func (s *testingStorage) Remove(migration *Migration) error { return nil }

type testingExecutor struct{}

func (e *testingExecutor) Up(*Migration, Storage) error   { return nil }
func (e *testingExecutor) Down(*Migration, Storage) error { return nil }

type stubbedExecutor struct {
	up   func(*Migration, Storage) error
	down func(*Migration, Storage) error
}

func (e *stubbedExecutor) Up(m *Migration, s Storage) error {
	if e.up != nil {
		return e.up(m, s)
	}

	return nil
}

func (e *stubbedExecutor) Down(m *Migration, s Storage) error {
	if e.down != nil {
		e.down(m, s)
	}

	return nil
}

func TestUnapplied(t *testing.T) {
	gl.Storage = &testingStorage{applied: Migrations{}}

	migrations, err := gl.Unapplied()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if migrations[0].Version != 20170329154959 {
		t.Fatalf("Expected version 20170329154959, got: %d", migrations[0].Version)
	}
}

func TestUnapplied_Empty(t *testing.T) {
	gl.Storage = &testingStorage{
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
	gl.Storage = &testingStorage{
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
	gl.Storage = &testingStorage{}

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

	gl.Storage = &testingStorage{}
	gl.Executor = &stubbedExecutor{
		up: func(*Migration, Storage) error {
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

	gl.Storage = &testingStorage{}
	gl.Executor = &stubbedExecutor{
		down: func(*Migration, Storage) error {
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
		InitialPath: "testdata/migrations",
		Source:      NewFileSystemSource("testdata/migrations"),
		Executor:    &testingExecutor{},
	}
}
