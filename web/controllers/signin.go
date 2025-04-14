package controllers

import (
	"net/http"
	"strings"

	"github.com/ordinary-dev/phoenix/web/sessions"
)

func (c *Controllers) ShowSignInForm(w http.ResponseWriter, _ *http.Request) {
	c.render("auth.html.tmpl", w, map[string]any{
		"title":       "Sign in",
		"description": "Authorization is required to view this page.",
		"button":      "Sign in",
		"formAction":  "/signin",
	})
}

func (c *Controllers) AuthorizeUser(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimSpace(r.FormValue("username"))
	password := strings.TrimSpace(r.FormValue("password"))
	user, err := c.db.GetUserIfPasswordMatches(username, password)
	if err != nil {
		c.ShowError(w, http.StatusUnauthorized, err)
		return
	}

	session, err := c.db.CreateSession(user.Username)
	if err != nil {
		c.ShowError(w, http.StatusInternalServerError, err)
		return
	}
	http.SetCookie(w, sessions.SessionToCookie(session))

	http.Redirect(w, r, "/", http.StatusFound)
}
