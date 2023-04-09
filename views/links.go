package views

import (
	"github.com/gin-gonic/gin"
	"github.com/ordinary-dev/phoenix/backend"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func CreateLink(c *gin.Context, db *gorm.DB) {
	if err := RequireAuth(c, db); err != nil {
		return
	}

	linkName := c.PostForm("linkName")
	href := c.PostForm("href")
	groupID, err := strconv.ParseUint(c.PostForm("groupID"), 10, 32)
	if err != nil {
		ShowError(c, err)
		return
	}

	if _, err = backend.CreateLink(db, linkName, href, groupID); err != nil {
		ShowError(c, err)
		return
	}

	// Redirect to settings.
	c.Redirect(http.StatusFound, "/settings")
}

func UpdateLink(c *gin.Context, db *gorm.DB) {
	if err := RequireAuth(c, db); err != nil {
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		ShowError(c, err)
		return
	}
	linkName := c.PostForm("linkName")
	href := c.PostForm("href")

	if _, err = backend.UpdateLink(db, id, linkName, href); err != nil {
		ShowError(c, err)
		return
	}

	// Redirect to settings.
	c.Redirect(http.StatusFound, "/settings")
}

func DeleteLink(c *gin.Context, db *gorm.DB) {
	if err := RequireAuth(c, db); err != nil {
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		ShowError(c, err)
		return
	}

	if err = backend.DeleteLink(db, id); err != nil {
		ShowError(c, err)
		return
	}

	// Redirect to settings.
	c.Redirect(http.StatusFound, "/settings")
}
