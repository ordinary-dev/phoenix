package database

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestMigrations(t *testing.T) {
	initTestDatabase(t)
	defer deleteTestDatabase(t)

	// We should be able to call the function multiple times.
	if err := ApplyMigrations(); err != nil {
		t.Fatal(err)
	}
}
