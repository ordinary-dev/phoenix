package views

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ordinary-dev/phoenix/backend"
	"gorm.io/gorm"
	"net/http"
)

func ShowRegistrationForm(c *gin.Context) {
	c.HTML(http.StatusOK, "auth.html.tmpl", gin.H{
		"title":       "Create an account",
		"description": "To prevent other people from seeing your links, create an account.",
		"button":      "Create",
		"formAction":  "/users",
	})
}

func ShowLoginForm(c *gin.Context) {
	c.HTML(http.StatusOK, "auth.html.tmpl", gin.H{
		"title":       "Sign in",
		"description": "Authorization is required to view this page.",
		"button":      "Sign in",
		"formAction":  "/signin",
	})
}

// Requires the user to log in before viewing the page.
// In case of an error, it shows the login page or the error page.
// Returns error if the user is not authorized.
// If `nil` is returned instead of an error, it is safe to display protected content.
func RequireAuth(c *gin.Context, db *gorm.DB) error {
	number_of_accounts := backend.CountAdmins(db)

	// First run
	if number_of_accounts == 0 {
		ShowRegistrationForm(c)
	}

	tokenValue, err := c.Cookie("phoenix-token")

	// Anonymous visitor
	if err != nil {
		ShowLoginForm(c)
		return errors.New("User is not authorized")
	}

	err = backend.ValidateToken(db, tokenValue)
	if err != nil {
		ShowError(c, err)
		return errors.New("Access token is invalid")
	}

	return nil
}
