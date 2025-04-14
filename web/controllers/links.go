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

func (c *Controllers) CreateLink(w http.ResponseWriter, r *http.Request) {
	groupID, err := strconv.Atoi(r.FormValue("groupID"))
	if err != nil {
		c.ShowError(w, http.StatusBadRequest, err)
		return
	}

	group, err := c.db.GetGroup(int(groupID))
	if err != nil {
		c.ShowError(w, http.StatusInternalServerError, err)
		return
	}

	username := sessions.GetUsername(r.Context())
	if group.Username == nil || *group.Username != username {
		c.ShowError(w, http.StatusForbidden, errors.New("you are not the owner of the group"))
		return
	}

	link := entities.Link{
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
	if err := c.db.CreateLink(&link); err != nil {
		c.ShowError(w, http.StatusInternalServerError, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/settings#link-%v", link.ID), http.StatusFound)
}

func (c *Controllers) UpdateLink(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		c.ShowError(w, http.StatusBadRequest, err)
		return
	}

	link, err := c.db.GetLink(id)
	if err != nil {
		c.ShowError(w, http.StatusInternalServerError, err)
		return
	}

	group, err := c.db.GetGroup(int(link.GroupID))
	if err != nil {
		c.ShowError(w, http.StatusInternalServerError, err)
		return
	}

	username := sessions.GetUsername(r.Context())
	if group.Username == nil || *group.Username != username {
		c.ShowError(w, http.StatusForbidden, errors.New("you are not the owner of the group"))
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

	if err := c.db.UpdateLink(link); err != nil {
		c.ShowError(w, http.StatusInternalServerError, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/settings#link-%v", link.ID), http.StatusFound)
}

func (c *Controllers) DeleteLink(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		c.ShowError(w, http.StatusBadRequest, err)
		return
	}

	link, err := c.db.GetLink(id)
	if err != nil {
		c.ShowError(w, http.StatusInternalServerError, err)
		return
	}

	group, err := c.db.GetGroup(int(link.GroupID))
	if err != nil {
		c.ShowError(w, http.StatusInternalServerError, err)
		return
	}

	username := sessions.GetUsername(r.Context())
	if group.Username == nil || *group.Username != username {
		c.ShowError(w, http.StatusForbidden, errors.New("you are not the owner of the group"))
		return
	}

	if err := c.db.DeleteLink(id); err != nil {
		c.ShowError(w, http.StatusInternalServerError, err)
		return
	}

	http.Redirect(w, r, "/settings", http.StatusFound)
}
