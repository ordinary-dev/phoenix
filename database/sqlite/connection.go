package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ordinary-dev/phoenix/config"
	"github.com/ordinary-dev/phoenix/database"
)

// Sqlite-based implementation of the Database interface.
type SqliteDB struct {
	conn *sql.DB
}

var _ database.Database = (*SqliteDB)(nil) // verify interface compliance

func (db *SqliteDB) Connect(cfg *config.Config) error {
	var err error
	db.conn, err = sql.Open("sqlite3", cfg.DBPath)
	if err != nil {
		return err
	}

	return db.conn.Ping()
}
