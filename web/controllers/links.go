package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ordinary-dev/phoenix/database"
)

func CreateLink(w http.ResponseWriter, r *http.Request) {
	groupID, err := strconv.Atoi(r.FormValue("groupID"))
	if err != nil {
		ShowError(w, http.StatusBadRequest, err)
		return
	}

	link := database.Link{
		Name:    strings.TrimSpace(r.FormValue("linkName")),
		Href:    strings.TrimSpace(r.FormValue("href")),
		GroupID: groupID,
	}
	icon := strings.TrimSpace(r.FormValue("icon"))
	if icon == "" {
		link.Icon = nil
	} else {
		link.Icon = &icon
	}
	if err := database.CreateLink(&link); err != nil {
		ShowError(w, http.StatusInternalServerError, err)
		return
	}

	// Redirect to settings.
	http.Redirect(w, r, fmt.Sprintf("/settings#link-%v", link.ID), http.StatusFound)
}

func UpdateLink(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		ShowError(w, http.StatusBadRequest, err)
		return
	}

	link, err := database.GetLink(id)
	if err != nil {
		ShowError(w, http.StatusInternalServerError, err)
		return
	}

	link.Name = strings.TrimSpace(r.FormValue("linkName"))
	link.Href = strings.TrimSpace(r.FormValue("href"))
	icon := strings.TrimSpace(r.FormValue("icon"))
	if icon == "" {
		link.Icon = nil
	} else {
		link.Icon = &icon
	}

	if err := database.UpdateLink(link); err != nil {
		ShowError(w, http.StatusInternalServerError, err)
		return
	}

	// Redirect to settings.
	http.Redirect(w, r, fmt.Sprintf("/settings#link-%v", link.ID), http.StatusFound)
}

func DeleteLink(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		ShowError(w, http.StatusBadRequest, err)
		return
	}

	if err := database.DeleteLink(id); err != nil {
		ShowError(w, http.StatusInternalServerError, err)
		return
	}

	// Redirect to settings.
	http.Redirect(w, r, "/settings", http.StatusFound)
}
