package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/ordinary-dev/phoenix/database/entities"
	"github.com/ordinary-dev/phoenix/web/sessions"
)

type ExportFile struct {
	Groups []entities.Group `json:"groups"`
}

func (c *Controllers) Export(w http.ResponseWriter, r *http.Request) {
	username := sessions.GetUsername(r.Context())
	groups, err := c.db.GetGroupsWithLinks(&username)
	if err != nil {
		c.ShowError(w, http.StatusInternalServerError, err)
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
