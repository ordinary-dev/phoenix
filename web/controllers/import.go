package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/ordinary-dev/phoenix/web/sessions"
)

func (c *Controllers) ImportPage(w http.ResponseWriter, _ *http.Request) {
	c.render("import.html.tmpl", w, map[string]any{})
}

func (c *Controllers) Import(w http.ResponseWriter, r *http.Request) {
	var exportFile ExportFile
	data := []byte(r.FormValue("exportFile"))
	if err := json.Unmarshal(data, &exportFile); err != nil {
		c.ShowError(w, http.StatusBadRequest, err)
		return
	}

	username := sessions.GetUsername(r.Context())
	for _, g := range exportFile.Groups {
		g.Username = &username
		if err := c.db.CreateGroup(&g); err != nil {
			c.ShowError(w, http.StatusInternalServerError, err)
			return
		}

		for _, l := range g.Links {
			l.GroupID = g.ID
			if err := c.db.CreateLink(&l); err != nil {
				c.ShowError(w, http.StatusInternalServerError, err)
				return
			}
		}
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
