package views

import (
	"github.com/gin-gonic/gin"
	"github.com/ordinary-dev/phoenix/backend"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func CreateGroup(c *gin.Context, db *gorm.DB) {
	if err := RequireAuth(c, db); err != nil {
		return
	}

	// Save new group to the database.
	groupName := c.PostForm("groupName")
	if _, err := backend.CreateGroup(db, groupName); err != nil {
		ShowError(c, err)
		return
	}

	// This page is called from the settings, return the user back.
	c.Redirect(http.StatusFound, "/settings")
}

func UpdateGroup(c *gin.Context, db *gorm.DB) {
	if err := RequireAuth(c, db); err != nil {
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		ShowError(c, err)
		return
	}
	groupName := c.PostForm("groupName")
	if _, err := backend.UpdateGroup(db, id, groupName); err != nil {
		ShowError(c, err)
		return
	}

	// This page is called from the settings, return the user back.
	c.Redirect(http.StatusFound, "/settings")
}

func DeleteGroup(c *gin.Context, db *gorm.DB) {
	if err := RequireAuth(c, db); err != nil {
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		ShowError(c, err)
		return
	}

	if err = backend.DeleteGroup(db, id); err != nil {
		ShowError(c, err)
		return
	}

	// Redirect to settings.
	c.Redirect(http.StatusFound, "/settings")
}
