package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"testing"
)

const TEST_DB_PATH = "/tmp/phoenix.sqlite3"

func initTestDatabase(t *testing.T) {
	var err error
	DB, err = sql.Open("sqlite3", TEST_DB_PATH)
	if err != nil {
		t.Fatal(err)
	}

	if err := ApplyMigrations(); err != nil {
		t.Fatal(err)
	}
}

func deleteTestDatabase(t *testing.T) {
	if err := os.Remove(TEST_DB_PATH); err != nil {
		t.Fatal(err)
	}
}
