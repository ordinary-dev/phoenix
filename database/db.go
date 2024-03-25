package database

import (
	"github.com/ordinary-dev/phoenix/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func EstablishDatabaseConnection(cfg *config.Config) error {
	var err error
	DB, err = gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{})
	if err != nil {
		return err
	}

	// Migrate the schema
	DB.AutoMigrate(&Admin{}, &Group{}, &Link{})

	return nil
}
