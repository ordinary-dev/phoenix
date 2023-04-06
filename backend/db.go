package backend

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func GetDatabaseConnection() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Migrate the schema
	db.AutoMigrate(&Admin{}, &AccessToken{}, &Group{}, &Link{})

	return db, nil
}
