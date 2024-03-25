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

	if err := database.CreateGroup(&group); err != nil {
		ShowError(w, http.StatusInternalServerError, err)
		return
	}

	// This page is called from the settings, return the user back.
	http.Redirect(w, r, fmt.Sprintf("/settings#group-%v", group.ID), http.StatusFound)
}

func UpdateGroup(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		ShowError(w, http.StatusBadRequest, err)
		return
	}

	if err := database.UpdateGroup(int(id), r.FormValue("groupName")); err != nil {
		ShowError(w, http.StatusInternalServerError, err)
		return
	}

	// This page is called from the settings, return the user back.
	http.Redirect(w, r, fmt.Sprintf("/settings#group-%v", id), http.StatusFound)
}

func DeleteGroup(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		ShowError(w, http.StatusBadRequest, err)
		return
	}

	if err := database.DeleteGroup(int(id)); err != nil {
		ShowError(w, http.StatusInternalServerError, err)
		return
	}

	// Redirect to settings.
	http.Redirect(w, r, "/settings", http.StatusFound)
}
