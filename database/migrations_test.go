package database

import (
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

func TestMigrations(t *testing.T) {
	initTestDatabase(t)
	defer deleteTestDatabase(t)

	// We should be able to call the function multiple times.
	if err := ApplyMigrations(); err != nil {
		t.Fatal(err)
	}
}
