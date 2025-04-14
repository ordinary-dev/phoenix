package controllers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/ordinary-dev/phoenix/config"
	"github.com/ordinary-dev/phoenix/web/sessions"
)

func (c *Controllers) ShowRegistrationForm(w http.ResponseWriter, _ *http.Request) {
	if !config.Cfg.EnableRegistration {
		c.ShowError(w, http.StatusForbidden, errors.New("registration disabled"))
		return
	}

	c.render("auth.html.tmpl", w, map[string]any{
		"title":       "Create an account",
		"description": "To prevent other people from seeing your links, create an account.",
		"button":      "Create",
		"formAction":  "/registration",
	})
}

func (c *Controllers) CreateUser(w http.ResponseWriter, r *http.Request) {
	if !config.Cfg.EnableRegistration {
		c.ShowError(w, http.StatusForbidden, errors.New("registration disabled"))
		return
	}

	username := strings.TrimSpace(r.FormValue("username"))
	password := strings.TrimSpace(r.FormValue("password"))
	user, err := c.db.CreateUser(username, &password)
	if err != nil {
		c.ShowError(w, http.StatusInternalServerError, err)
		return
	}

	// Generate access token.
	session, err := c.db.CreateSession(user.Username)
	if err != nil {
		c.ShowError(w, http.StatusInternalServerError, err)
		return
	}
	http.SetCookie(w, sessions.SessionToCookie(session))

	http.Redirect(w, r, "/", http.StatusFound)
}
