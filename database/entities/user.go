package entities

import (
	"time"
)

type User struct {
	Username       string
	HashedPassword *string
}

type Session struct {
	Token     string
	Username  string
	CreatedAt time.Time
}

const (
	TokenLength   = 32
	TokenLifetime = time.Hour * 24 * 30
)
