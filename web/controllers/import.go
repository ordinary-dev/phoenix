package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/ordinary-dev/phoenix/database"
	"github.com/ordinary-dev/phoenix/web/sessions"
)

func ImportPage(w http.ResponseWriter, _ *http.Request) {
	Render("import.html.tmpl", w, map[string]any{})
}

func Import(w http.ResponseWriter, r *http.Request) {
	var exportFile ExportFile
	data := []byte(r.FormValue("exportFile"))
	if err := json.Unmarshal(data, &exportFile); err != nil {
		ShowError(w, http.StatusBadRequest, err)
		return
	}

	username := sessions.GetUsername(r.Context())
	for _, g := range exportFile.Groups {
		g.Username = &username
		if err := database.CreateGroup(&g); err != nil {
			ShowError(w, http.StatusInternalServerError, err)
			return
		}

		for _, l := range g.Links {
			l.GroupID = g.ID
			if err := database.CreateLink(&l); err != nil {
				ShowError(w, http.StatusInternalServerError, err)
				return
			}
		}
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
