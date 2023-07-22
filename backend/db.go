package backend

import (
	"github.com/ordinary-dev/phoenix/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func GetDatabaseConnection(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Migrate the schema
	db.AutoMigrate(&Admin{}, &Group{}, &Link{})

	return db, nil
}
