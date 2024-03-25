package pages

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ordinary-dev/phoenix/database"
)

func CreateGroup(w http.ResponseWriter, r *http.Request) {
	// Save new group to the database.
	group := database.Group{
		Name: r.FormValue("groupName"),
	}

	if result := database.DB.Create(&group); result.Error != nil {
		ShowError(w, http.StatusInternalServerError, result.Error)
		return
	}

	// This page is called from the settings, return the user back.
	http.Redirect(w, r, fmt.Sprintf("/settings#group-%v", group.ID), http.StatusFound)
}

func UpdateGroup(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		ShowError(w, http.StatusBadRequest, err)
		return
	}

	var group database.Group
	if result := database.DB.First(&group, id); result.Error != nil {
		ShowError(w, http.StatusInternalServerError, result.Error)
		return
	}

	group.Name = r.FormValue("groupName")
	if result := database.DB.Save(&group); result.Error != nil {
		ShowError(w, http.StatusInternalServerError, result.Error)
		return
	}

	// This page is called from the settings, return the user back.
	http.Redirect(w, r, fmt.Sprintf("/settings#group-%v", group.ID), http.StatusFound)
}

func DeleteGroup(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		ShowError(w, http.StatusBadRequest, err)
		return
	}

	if result := database.DB.Delete(&database.Group{}, id); result.Error != nil {
		ShowError(w, http.StatusInternalServerError, result.Error)
		return
	}

	// Redirect to settings.
	http.Redirect(w, r, "/settings", http.StatusFound)
}
