package main

import "os"

const helpMsg = `Usage gloat: [COMMAND ...]

Gloat is a Go SQL migration utility.

Commands:
  up            Apply new migrations
  down          Revert the last applied migration

Options:
  --help        Show this message
`

func main() {
	args := os.Args

	if len(args) < 2 {
		Exitf(1, helpMsg)
	}

	switch args[1] {
	case "up":
		upCmd()
	case "down":
		downCmd()
	}
}

func upCmd() {
	migrations, err := gl.UnappliedMigrations()
	if err != nil {
		Exitf(1, "Error: %v\n", err)
	}

	appliedMigrations := map[int64]bool{}

	for _, migration := range migrations {
		Outf("Applying migration: %d...\n", migration.Version)

		if err := gl.Apply(migration); err != nil {
			Exitf(1, "Error: %v\n", err)
		}

		appliedMigrations[migration.Version] = true
	}

	if len(appliedMigrations) == 0 {
		Outf("No migrations to apply\n")
	}
}

func downCmd() {
	migration, err := gl.CurrentMigration()
	if err != nil {
		Exitf(1, "Error: %v\n", err)
	}

	if migration == nil {
		Exitf(0, "No migrations to apply\n")
	}

	Outf("Reverting migration: %d...\n", migration.Version)

	if err := gl.Revert(migration); err != nil {
		Exitf(1, "Error: %v\n", err)
	}
}
