package views

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ordinary-dev/phoenix/config"
	"github.com/ordinary-dev/phoenix/database"
	"gorm.io/gorm"
)

const TOKEN_LIFETIME_IN_SECONDS = 60 * 60 * 24 * 30

func ShowRegistrationForm(cfg *config.Config, db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if database.CountAdmins(db) > 0 {
			ShowError(ctx, cfg, errors.New("At least 1 user already exists"))
			return
		}

		Render(ctx, cfg, http.StatusOK, "auth.html.tmpl", gin.H{
			"title":       "Create an account",
			"description": "To prevent other people from seeing your links, create an account.",
			"button":      "Create",
			"formAction":  "/api/users",
		})
	}
}

func ShowLoginForm(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		Render(ctx, cfg, http.StatusOK, "auth.html.tmpl", gin.H{
			"title":       "Sign in",
			"description": "Authorization is required to view this page.",
			"button":      "Sign in",
			"formAction":  "/api/users/signin",
		})
	}
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

func AuthMiddleware(db *gorm.DB, cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, err := RequireAuth(ctx, cfg)
		if err != nil {
			if cfg.HeaderAuth && ctx.Request.Header.Get("Remote-User") != "" {
				// Generate access token.
				token, err := GetJWTToken(cfg)
				if err != nil {
					ShowError(ctx, cfg, err)
					return
				}
				SetTokenCookie(ctx, token, cfg)
				return
			}

			if database.CountAdmins(db) < 1 {
				ctx.Redirect(http.StatusFound, "/registration")
			} else {
				ctx.Redirect(http.StatusFound, "/signin")
			}
			ctx.Abort()
			return
		}

		// Create a new token if the old one is about to expire
		if time.Now().Add(time.Second * (TOKEN_LIFETIME_IN_SECONDS / 2)).After(claims.ExpiresAt.Time) {
			newToken, err := GetJWTToken(cfg)
			if err != nil {
				ShowError(ctx, cfg, err)
				return
			}
			SetTokenCookie(ctx, newToken, cfg)
		}
	}
}

func GetJWTToken(cfg *config.Config) (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * TOKEN_LIFETIME_IN_SECONDS)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.SecretKey))
}

func CreateUser(db *gorm.DB, cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if database.CountAdmins(db) > 0 {
			ShowError(ctx, cfg, errors.New("At least 1 user already exists"))
			return
		}

		// Try to create a user.
		username := ctx.PostForm("username")
		password := ctx.PostForm("password")
		_, err := database.CreateAdmin(db, username, password)
		if err != nil {
			ShowError(ctx, cfg, err)
			return
		}

		// Generate access token.
		token, err := GetJWTToken(cfg)
		if err != nil {
			ShowError(ctx, cfg, err)
			return
		}
		SetTokenCookie(ctx, token, cfg)

		// Redirect to homepage.
		ctx.Redirect(http.StatusFound, "/")
	}
}

func AuthorizeUser(db *gorm.DB, cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Check credentials.
		username := ctx.PostForm("username")
		password := ctx.PostForm("password")
		_, err := database.AuthorizeAdmin(db, username, password)
		if err != nil {
			ShowError(ctx, cfg, err)
			return
		}

		// Generate an access token.
		token, err := GetJWTToken(cfg)
		if err != nil {
			ShowError(ctx, cfg, err)
			return
		}
		SetTokenCookie(ctx, token, cfg)

		// Redirect to homepage.
		ctx.Redirect(http.StatusFound, "/")
	}
}

// Save token in cookies
func SetTokenCookie(c *gin.Context, token string, cfg *config.Config) {
	c.SetCookie("phoenix-token", token, TOKEN_LIFETIME_IN_SECONDS, "/", "", cfg.SecureCookie, true)
}
