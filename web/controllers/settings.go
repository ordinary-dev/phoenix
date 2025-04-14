package controllers

import (
	"net/http"

	"github.com/ordinary-dev/phoenix/web/sessions"
)

func (c *Controllers) ShowSettings(w http.ResponseWriter, r *http.Request) {
	username := sessions.GetUsername(r.Context())
	groups, err := c.db.GetGroupsWithLinks(&username)
	if err != nil {
		c.ShowError(w, http.StatusInternalServerError, err)
		return
	}

	c.render("settings.html.tmpl", w, map[string]any{
		"title":  "Settings",
		"groups": groups,
	})
}
