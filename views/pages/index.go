package pages

import (
	"net/http"

	"github.com/ordinary-dev/phoenix/database"
	log "github.com/sirupsen/logrus"
)

func ShowMainPage(w http.ResponseWriter, _ *http.Request) {
	groups, err := database.GetGroupsWithLinks()
	if err != nil {
		ShowError(w, http.StatusInternalServerError, err)
		return
	}

	err = Render("index.html.tmpl", w, map[string]any{
		"description": "Self-hosted start page.",
		"groups":      groups,
	})
	if err != nil {
		log.Error(err)
	}
}
