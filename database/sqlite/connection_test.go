package sqlite

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

const TEST_DB_PATH = "/tmp/phoenix.sqlite3"

func initTestDatabase(t *testing.T) SqliteDB {
	var err error
	var db SqliteDB
	db.conn, err = sql.Open("sqlite3", TEST_DB_PATH)
	if err != nil {
		t.Fatal(err)
	}

	if err := db.Migrate(); err != nil {
		t.Fatal(err)
	}

	return db
}

func deleteTestDatabase(t *testing.T) {
	if err := os.Remove(TEST_DB_PATH); err != nil {
		t.Fatal(err)
	}
}
