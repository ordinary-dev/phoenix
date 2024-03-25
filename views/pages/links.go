package pages

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ordinary-dev/phoenix/database"
)

func CreateLink(w http.ResponseWriter, r *http.Request) {
	groupID, err := strconv.ParseUint(r.FormValue("groupID"), 10, 32)
	if err != nil {
		ShowError(w, http.StatusBadRequest, err)
		return
	}

	link := database.Link{
		Name:    r.FormValue("linkName"),
		Href:    r.FormValue("href"),
		GroupID: groupID,
	}
	icon := r.FormValue("icon")
	if icon == "" {
		link.Icon = nil
	} else {
		link.Icon = &icon
	}
	if result := database.DB.Create(&link); result.Error != nil {
		ShowError(w, http.StatusInternalServerError, result.Error)
		return
	}

	// Redirect to settings.
	http.Redirect(w, r, fmt.Sprintf("/settings#link-%v", link.ID), http.StatusFound)
}

func UpdateLink(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		ShowError(w, http.StatusBadRequest, err)
		return
	}

	var link database.Link
	if result := database.DB.First(&link, id); result.Error != nil {
		ShowError(w, http.StatusInternalServerError, result.Error)
		return
	}

	link.Name = r.FormValue("linkName")
	link.Href = r.FormValue("href")
	icon := r.FormValue("icon")
	if icon == "" {
		link.Icon = nil
	} else {
		link.Icon = &icon
	}
	if result := database.DB.Save(&link); result.Error != nil {
		ShowError(w, http.StatusInternalServerError, result.Error)
		return
	}

	// Redirect to settings.
	http.Redirect(w, r, fmt.Sprintf("/settings#link-%v", link.ID), http.StatusFound)
}

func DeleteLink(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		ShowError(w, http.StatusBadRequest, err)
		return
	}

	if result := database.DB.Delete(&database.Link{}, id); result.Error != nil {
		ShowError(w, http.StatusInternalServerError, result.Error)
		return
	}

	// Redirect to settings.
	http.Redirect(w, r, "/settings", http.StatusFound)
}
