package backend

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Admin struct {
	ID       uint64 `gorm:"primaryKey"`
	Username string `gorm:"unique;notNull"`
	Bcrypt   string `gorm:"notNull"`
}

func CountAdmins(db *gorm.DB) int64 {
	var admins []Admin
	var count int64
	db.Model(&admins).Count(&count)
	return count
}

func CreateAdmin(db *gorm.DB, username string, password string) (Admin, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return Admin{}, err
	}

	admin := Admin{
		Username: username,
		Bcrypt:   string(hash),
	}
	result := db.Create(&admin)

	if result.Error != nil {
		return Admin{}, result.Error
	}

	return admin, nil
}

func AuthorizeAdmin(db *gorm.DB, username string, password string) (Admin, error) {
	var admin Admin
	result := db.Where("username = ?", username).First(&admin)

	if result.Error != nil {
		return Admin{}, result.Error
	}

	err := bcrypt.CompareHashAndPassword([]byte(admin.Bcrypt), []byte(password))
	if err != nil {
		return Admin{}, err
	}

	return admin, nil
}
