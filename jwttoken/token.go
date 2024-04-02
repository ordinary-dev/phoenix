package jwttoken

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/ordinary-dev/phoenix/config"
)

const (
	TOKEN_LIFETIME_IN_SECONDS = 60 * 60 * 24 * 30
	TOKEN_COOKIE_NAME         = "phoenix-token"
)

func GetJWTToken() (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * TOKEN_LIFETIME_IN_SECONDS)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Cfg.SecretKey))
}

func TokenToCookie(value string) *http.Cookie {
	return &http.Cookie{
		Name:     TOKEN_COOKIE_NAME,
		Value:    value,
		HttpOnly: true,
		Secure:   config.Cfg.SecureCookie,
		MaxAge:   TOKEN_LIFETIME_IN_SECONDS,
	}
}
