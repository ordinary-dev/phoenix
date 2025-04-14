package sqlite

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestMigrations(t *testing.T) {
	db := initTestDatabase(t)
	defer deleteTestDatabase(t)

	// We should be able to call the function multiple times.
	if err := db.Migrate(); err != nil {
		t.Fatal(err)
	}
}
