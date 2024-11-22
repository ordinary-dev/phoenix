package controllers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/ordinary-dev/phoenix/config"
	"github.com/ordinary-dev/phoenix/database"
	"github.com/ordinary-dev/phoenix/web/sessions"
)

func ShowRegistrationForm(w http.ResponseWriter, _ *http.Request) {
	if !config.Cfg.EnableRegistration {
		ShowError(w, http.StatusForbidden, errors.New("registration disabled"))
		return
	}

	Render("auth.html.tmpl", w, map[string]any{
		"title":       "Create an account",
		"description": "To prevent other people from seeing your links, create an account.",
		"button":      "Create",
		"formAction":  "/registration",
	})
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	if !config.Cfg.EnableRegistration {
		ShowError(w, http.StatusForbidden, errors.New("registration disabled"))
		return
	}

	username := strings.TrimSpace(r.FormValue("username"))
	password := strings.TrimSpace(r.FormValue("password"))
	user, err := database.CreateUser(username, &password)
	if err != nil {
		ShowError(w, http.StatusInternalServerError, err)
		return
	}

	// Generate access token.
	session, err := database.CreateSession(user.Username)
	if err != nil {
		ShowError(w, http.StatusInternalServerError, err)
		return
	}
	http.SetCookie(w, sessions.SessionToCookie(session))

	http.Redirect(w, r, "/", http.StatusFound)
}
