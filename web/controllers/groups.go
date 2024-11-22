package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ordinary-dev/phoenix/database"
	"github.com/ordinary-dev/phoenix/web/sessions"
)

func CreateGroup(w http.ResponseWriter, r *http.Request) {
	username := sessions.GetUsername(r.Context())
	group := database.Group{
		Name:     strings.TrimSpace(r.FormValue("groupName")),
		Username: &username,
	}

	if err := database.CreateGroup(&group); err != nil {
		ShowError(w, http.StatusInternalServerError, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/settings#group-%v", group.ID), http.StatusFound)
}

func UpdateGroup(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		ShowError(w, http.StatusBadRequest, err)
		return
	}

	group, err := database.GetGroup(int(id))
	if err != nil {
		ShowError(w, http.StatusInternalServerError, err)
		return
	}

	username := sessions.GetUsername(r.Context())
	if group.Username == nil || *group.Username != username {
		ShowError(w, http.StatusForbidden, errors.New("you are not the owner of the group"))
		return
	}

	newName := strings.TrimSpace(r.FormValue("groupName"))
	if err := database.UpdateGroup(int(id), newName); err != nil {
		ShowError(w, http.StatusInternalServerError, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/settings#group-%v", id), http.StatusFound)
}

func DeleteGroup(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		ShowError(w, http.StatusBadRequest, err)
		return
	}

	group, err := database.GetGroup(int(id))
	if err != nil {
		ShowError(w, http.StatusInternalServerError, err)
		return
	}

	username := sessions.GetUsername(r.Context())
	if group.Username == nil || *group.Username != username {
		ShowError(w, http.StatusForbidden, errors.New("you are not the owner of the group"))
		return
	}

	if err := database.DeleteGroup(int(id)); err != nil {
		ShowError(w, http.StatusInternalServerError, err)
		return
	}

	http.Redirect(w, r, "/settings", http.StatusFound)
}
