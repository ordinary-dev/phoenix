package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/ordinary-dev/phoenix/config"
	"github.com/ordinary-dev/phoenix/database"
	"github.com/ordinary-dev/phoenix/jwttoken"
	"github.com/ordinary-dev/phoenix/web/controllers"
)

// Try to find the access token in the request.
// Returns error if the user is not authorized.
// If `nil` is returned instead of an error, it is safe to display protected content.
func ParseToken(r *http.Request) (*jwt.RegisteredClaims, error) {
	tokenCookie, err := r.Cookie(jwttoken.TOKEN_COOKIE_NAME)

	// Anonymous visitor.
	if err != nil {
		return nil, err
	}

	// Check token.
	token, err := jwt.ParseWithClaims(tokenCookie.Value, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.Cfg.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	return claims, nil
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if SSO is enabled.
		if config.Cfg.HeaderAuth && r.Header.Get("Remote-User") != "" {
			next.ServeHTTP(w, r)
			return
		}

		claims, err := ParseToken(r)

		// Most likely the user is not authorized.
		if err != nil {
			count, err := database.CountAdmins()
			if err != nil {
				controllers.ShowError(w, http.StatusInternalServerError, err)
				return
			}

			if count < 1 {
				http.Redirect(w, r, "/registration", http.StatusFound)
			} else {
				http.Redirect(w, r, "/signin", http.StatusFound)
			}

			return
		}

		// Create a new token if the old one is about to expire
		if time.Now().Add(time.Second * (jwttoken.TOKEN_LIFETIME_IN_SECONDS / 2)).After(claims.ExpiresAt.Time) {
			newToken, err := jwttoken.GetJWTToken()
			if err != nil {
				controllers.ShowError(w, http.StatusInternalServerError, err)
				return
			}

			http.SetCookie(w, jwttoken.TokenToCookie(newToken))
		}

		next.ServeHTTP(w, r)
	})
}
