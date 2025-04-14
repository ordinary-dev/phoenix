package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ordinary-dev/phoenix/database/entities"
	"github.com/ordinary-dev/phoenix/web/sessions"
)

func (c *Controllers) CreateGroup(w http.ResponseWriter, r *http.Request) {
	username := sessions.GetUsername(r.Context())
	group := entities.Group{
		Name:     strings.TrimSpace(r.FormValue("groupName")),
		Username: &username,
	}

	if err := c.db.CreateGroup(&group); err != nil {
		c.ShowError(w, http.StatusInternalServerError, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/settings#group-%v", group.ID), http.StatusFound)
}

func (c *Controllers) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		c.ShowError(w, http.StatusBadRequest, err)
		return
	}

	group, err := c.db.GetGroup(int(id))
	if err != nil {
		c.ShowError(w, http.StatusInternalServerError, err)
		return
	}

	username := sessions.GetUsername(r.Context())
	if group.Username == nil || *group.Username != username {
		c.ShowError(w, http.StatusForbidden, errors.New("you are not the owner of the group"))
		return
	}

	newName := strings.TrimSpace(r.FormValue("groupName"))
	if err := c.db.UpdateGroup(int(id), newName); err != nil {
		c.ShowError(w, http.StatusInternalServerError, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/settings#group-%v", id), http.StatusFound)
}

func (c *Controllers) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		c.ShowError(w, http.StatusBadRequest, err)
		return
	}

	group, err := c.db.GetGroup(int(id))
	if err != nil {
		c.ShowError(w, http.StatusInternalServerError, err)
		return
	}

	username := sessions.GetUsername(r.Context())
	if group.Username == nil || *group.Username != username {
		c.ShowError(w, http.StatusForbidden, errors.New("you are not the owner of the group"))
		return
	}

	if err := c.db.DeleteGroup(int(id)); err != nil {
		c.ShowError(w, http.StatusInternalServerError, err)
		return
	}

	http.Redirect(w, r, "/settings", http.StatusFound)
}
