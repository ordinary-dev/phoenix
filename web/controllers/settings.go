package controllers

import (
	"net/http"

	"github.com/ordinary-dev/phoenix/database"
	"github.com/ordinary-dev/phoenix/web/sessions"
)

func ShowSettings(w http.ResponseWriter, r *http.Request) {
	username := sessions.GetUsername(r.Context())
	groups, err := database.GetGroupsWithLinks(&username)
	if err != nil {
		ShowError(w, http.StatusInternalServerError, err)
		return
	}

	Render("settings.html.tmpl", w, map[string]any{
		"title":  "Settings",
		"groups": groups,
	})
}
