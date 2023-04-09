package backend

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

func GetDatabaseConnection() (*gorm.DB, error) {
	dbPath := os.Getenv("PHOENIX_DB_PATH")
	if dbPath == "" {
		dbPath = "db.sqlite3"
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Migrate the schema
	db.AutoMigrate(&Admin{}, &AccessToken{}, &Group{}, &Link{})

	return db, nil
}
