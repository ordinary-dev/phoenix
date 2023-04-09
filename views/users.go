package views

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ordinary-dev/phoenix/backend"
	"gorm.io/gorm"
	"net/http"
)

func CreateUser(c *gin.Context, db *gorm.DB) {
	// If at least 1 administator exists, require authorization
	if backend.CountAdmins(db) > 0 {
		tokenValue, err := c.Cookie("phoenix-token")

		// Anonymous visitor
		if err != nil {
			err = errors.New("At least 1 user exists, you have to sign in first")
			ShowError(c, err)
			return
		}

		err = backend.ValidateToken(db, tokenValue)
		if err != nil {
			ShowError(c, err)
			return
		}
	}

	// User is authorized or no user exists.
	// Try to create a user.
	username := c.PostForm("username")
	password := c.PostForm("password")
	admin, err := backend.CreateAdmin(db, username, password)
	if err != nil {
		ShowError(c, err)
		return
	}

	// Generate access token.
	token, err := backend.CreateAccessToken(db, admin.ID)
	if err != nil {
		ShowError(c, err)
		return
	}
	SetTokenCookie(c, token)

	// Redirect to homepage.
	c.Redirect(http.StatusFound, "/")
}

func AuthorizeUser(c *gin.Context, db *gorm.DB) {
	// Check credentials.
	username := c.PostForm("username")
	password := c.PostForm("password")
	admin, err := backend.AuthorizeAdmin(db, username, password)
	if err != nil {
		ShowError(c, err)
		return
	}

	// Generate an access token.
	token, err := backend.CreateAccessToken(db, admin.ID)
	if err != nil {
		ShowError(c, err)
		return
	}
	SetTokenCookie(c, token)

	// Redirect to homepage.
	c.Redirect(http.StatusFound, "/")
}

// Save token for 29 days in cookies
func SetTokenCookie(c *gin.Context, token backend.AccessToken) {
	c.SetCookie("phoenix-token", token.Value, 60*60*24*29, "/", "", false, true)
}
