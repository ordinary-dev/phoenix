package views

import (
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
// If successful, does nothing.
// In case of an error, it shows the login page or the error page.
func RequireAuth(c *gin.Context, db *gorm.DB) {
	number_of_accounts := backend.CountAdmins(db)

	// First run
	if number_of_accounts == 0 {
		ShowRegistrationForm(c)
	}

	tokenValue, err := c.Cookie("phoenix-token")

	// Anonymous visitor
	if err != nil {
		ShowLoginForm(c)
		return
	}

	err = backend.ValidateToken(db, tokenValue)
	if err != nil {
		ShowError(c, err)
		return
	}
}
