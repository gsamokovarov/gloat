package gloat

import "testing"

var gl Gloat

type testingStorage struct{ applied Migrations }

func (s *testingStorage) Collect() (Migrations, error)      { return s.applied, nil }
func (s *testingStorage) Insert(migration *Migration) error { return nil }
func (s *testingStorage) Remove(migration *Migration) error { return nil }

type testingExecutor struct{}

func (s *testingExecutor) Up(*Migration, Storage) error   { return nil }
func (s *testingExecutor) Down(*Migration, Storage) error { return nil }

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

func TestUnappliedNone(t *testing.T) {
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

func init() {
	gl = Gloat{
		InitialPath: "testdata/migrations",
		Source:      NewFileSystemSource("testdata/migrations"),
		Executor:    &testingExecutor{},
	}
}
