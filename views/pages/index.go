package pages

import (
	"net/http"

	"github.com/ordinary-dev/phoenix/database"
	log "github.com/sirupsen/logrus"
)

func ShowMainPage(w http.ResponseWriter, _ *http.Request) {
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

	err := Render("index.html.tmpl", w, map[string]any{
		"description": "Self-hosted start page.",
		"groups":      groups,
	})
	if err != nil {
		log.Error(err)
	}
}
