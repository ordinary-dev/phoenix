package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ordinary-dev/phoenix/config"
)

var DB *sql.DB

func EstablishDatabaseConnection(cfg *config.Config) error {
	var err error
	DB, err = sql.Open("sqlite3", cfg.DBPath)

	if err := DB.Ping(); err != nil {
		return err
	}

	return err
}
