package views

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ordinary-dev/phoenix/backend"
	"github.com/ordinary-dev/phoenix/config"
	"gorm.io/gorm"
	"net/http"
	"time"
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
// Returns error if the user is not authorized.
// If `nil` is returned instead of an error, it is safe to display protected content.
func RequireAuth(c *gin.Context, cfg *config.Config) (*jwt.RegisteredClaims, error) {
	tokenValue, err := c.Cookie("phoenix-token")

	// Anonymous visitor
	if err != nil {
		return nil, err
	}

	// Check token
	token, err := jwt.ParseWithClaims(tokenValue, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(cfg.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return nil, errors.New("Token is invalid")
	}

	return claims, nil
}

func AuthMiddleware(c *gin.Context, cfg *config.Config) {
	claims, err := RequireAuth(c, cfg)
	if err != nil {
		c.Redirect(http.StatusFound, "/signin")
		c.Abort()
		return
	}

	// Create a new token if the old one is about to expire
	if time.Now().Add(12 * time.Hour).After(claims.ExpiresAt.Time) {
		newToken, err := GetJWTToken(cfg)
		if err != nil {
			ShowError(c, err)
			return
		}
		SetTokenCookie(c, newToken)
	}
}

func GetJWTToken(cfg *config.Config) (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.SecretKey))
}

func CreateUser(c *gin.Context, db *gorm.DB, cfg *config.Config) {
	// Try to create a user.
	username := c.PostForm("username")
	password := c.PostForm("password")
	_, err := backend.CreateAdmin(db, username, password)
	if err != nil {
		ShowError(c, err)
		return
	}

	// Generate access token.
	token, err := GetJWTToken(cfg)
	if err != nil {
		ShowError(c, err)
		return
	}
	SetTokenCookie(c, token)

	// Redirect to homepage.
	c.Redirect(http.StatusFound, "/")
}

func AuthorizeUser(c *gin.Context, db *gorm.DB, cfg *config.Config) {
	// Check credentials.
	username := c.PostForm("username")
	password := c.PostForm("password")
	_, err := backend.AuthorizeAdmin(db, username, password)
	if err != nil {
		ShowError(c, err)
		return
	}

	// Generate an access token.
	token, err := GetJWTToken(cfg)
	if err != nil {
		ShowError(c, err)
		return
	}
	SetTokenCookie(c, token)

	// Redirect to homepage.
	c.Redirect(http.StatusFound, "/")
}

// Save token for 29 days in cookies
func SetTokenCookie(c *gin.Context, token string) {
	c.SetCookie("phoenix-token", token, 60*60*24*29, "/", "", false, true)
}
