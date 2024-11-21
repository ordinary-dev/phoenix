package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/ordinary-dev/phoenix/database"
	"github.com/ordinary-dev/phoenix/web/sessions"
)

type ExportFile struct {
	Groups []database.Group `json:"groups"`
}

func Export(w http.ResponseWriter, r *http.Request) {
	username := sessions.GetUsername(r.Context())
	groups, err := database.GetGroupsWithLinks(&username)
	if err != nil {
		ShowError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename=phoenix.json")
	w.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(w)
	enc.Encode(&ExportFile{
		Groups: groups,
	})
}
