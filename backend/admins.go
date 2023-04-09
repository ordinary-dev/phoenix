package backend

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type Admin struct {
	ID           uint64 `gorm:"primaryKey"`
	Username     string `gorm:"unique;notNull"`
	Bcrypt       string `gorm:"notNull"`
	AccessTokens []AccessToken
}

type AccessToken struct {
	ID         uint64    `gorm:"primaryKey"`
	Value      string    `gorm:"notNull"`
	AdminID    uint64    `gorm:"notNull"`
	ValidUntil time.Time `gorm:"NotNull"`
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

func CreateAccessToken(db *gorm.DB, adminID uint64) (AccessToken, error) {
	bytes := make([]byte, 64)
	if _, err := rand.Read(bytes); err != nil {
		return AccessToken{}, err
	}

	accessToken := AccessToken{
		AdminID: adminID,
		Value:   base64.StdEncoding.EncodeToString(bytes),
		// Valid for 1 month
		ValidUntil: time.Now().AddDate(0, 1, 0),
	}
	result := db.Create(&accessToken)

	if result.Error != nil {
		return AccessToken{}, result.Error
	}

	return accessToken, nil
}

func ValidateToken(db *gorm.DB, value string) error {
	var token AccessToken
	result := db.Where("value = ?", value).First(&token)

	if result.Error != nil {
		return result.Error
	}

	if time.Now().After(token.ValidUntil) {
		return errors.New("Access token expired")
	}

	return nil
}
