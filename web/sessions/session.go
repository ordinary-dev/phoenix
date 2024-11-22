package sessions

import (
	"net/http"

	"github.com/ordinary-dev/phoenix/config"
	"github.com/ordinary-dev/phoenix/database"
)

const (
	TokenCookieName = "phoenix-session"
)

func SessionToCookie(session database.Session) *http.Cookie {
	return &http.Cookie{
		Name:     TokenCookieName,
		Value:    session.Token,
		Secure:   config.Cfg.SecureCookie,
		HttpOnly: true,
		Expires:  session.CreatedAt.Add(database.TokenLifetime),
	}
}
