package controllers

import (
	"net/http"
	"strings"

	"github.com/ordinary-dev/phoenix/database"
	"github.com/ordinary-dev/phoenix/web/sessions"
)

func ShowSignInForm(w http.ResponseWriter, _ *http.Request) {
	Render("auth.html.tmpl", w, map[string]any{
		"title":       "Sign in",
		"description": "Authorization is required to view this page.",
		"button":      "Sign in",
		"formAction":  "/signin",
	})
}

func AuthorizeUser(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimSpace(r.FormValue("username"))
	password := strings.TrimSpace(r.FormValue("password"))
	user, err := database.GetUserIfPasswordMatches(username, password)
	if err != nil {
		ShowError(w, http.StatusUnauthorized, err)
		return
	}

	session, err := database.CreateSession(user.Username)
	if err != nil {
		ShowError(w, http.StatusInternalServerError, err)
		return
	}
	http.SetCookie(w, sessions.SessionToCookie(session))

	http.Redirect(w, r, "/", http.StatusFound)
}
