package pages

import (
	"net/http"

	"github.com/ordinary-dev/phoenix/database"
)

func ShowSettings(w http.ResponseWriter, _ *http.Request) {
	// Get a list of groups with links
	var groups []database.Group
	result := database.DB.
		Model(&database.Group{}).
		Preload("Links").
		Find(&groups)

	if result.Error != nil {
		ShowError(w, http.StatusInternalServerError, result.Error)
		return
	}

	Render("settings.html.tmpl", w, map[string]any{
		"title":  "Settings",
		"groups": groups,
	})
}
