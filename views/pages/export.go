package pages

import (
	"encoding/json"
	"net/http"

	"github.com/ordinary-dev/phoenix/database"
)

type ExportFile struct {
	Groups []database.Group `json:"groups"`
}

func Export(w http.ResponseWriter, _ *http.Request) {
	groups, err := database.GetGroupsWithLinks()
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
