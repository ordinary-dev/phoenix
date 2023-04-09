package views

import (
	"github.com/gin-gonic/gin"
	"github.com/ordinary-dev/phoenix/backend"
	"gorm.io/gorm"
	"net/http"
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
