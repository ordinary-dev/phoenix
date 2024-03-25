package pages

import (
	"net/http"

	"github.com/ordinary-dev/phoenix/database"
)

func ShowSettings(w http.ResponseWriter, _ *http.Request) {
	groups, err := database.GetGroupsWithLinks()
	if err != nil {
		ShowError(w, http.StatusInternalServerError, err)
		return
	}

	Render("settings.html.tmpl", w, map[string]any{
		"title":  "Settings",
		"groups": groups,
	})
}
