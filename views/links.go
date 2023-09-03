package views

import (
	"github.com/gin-gonic/gin"
	"github.com/ordinary-dev/phoenix/database"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func CreateLink(c *gin.Context, db *gorm.DB) {
	groupID, err := strconv.ParseUint(c.PostForm("groupID"), 10, 32)
	if err != nil {
		ShowError(c, err)
		return
	}

	link := database.Link{
		Name:    c.PostForm("linkName"),
		Href:    c.PostForm("href"),
		GroupID: groupID,
	}
	icon := c.PostForm("icon")
	if icon == "" {
		link.Icon = nil
	} else {
		link.Icon = &icon
	}
	if result := db.Create(&link); result.Error != nil {
		ShowError(c, result.Error)
		return
	}

	// Redirect to settings.
	c.Redirect(http.StatusFound, "/settings")
}

func UpdateLink(c *gin.Context, db *gorm.DB) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		ShowError(c, err)
		return
	}

	var link database.Link
	if result := db.First(&link, id); result.Error != nil {
		ShowError(c, err)
		return
	}

	link.Name = c.PostForm("linkName")
	link.Href = c.PostForm("href")
	icon := c.PostForm("icon")
	if icon == "" {
		link.Icon = nil
	} else {
		link.Icon = &icon
	}
	if result := db.Save(&link); result.Error != nil {
		ShowError(c, result.Error)
		return
	}

	// Redirect to settings.
	c.Redirect(http.StatusFound, "/settings")
}

func DeleteLink(c *gin.Context, db *gorm.DB) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		ShowError(c, err)
		return
	}

	if result := db.Delete(&database.Link{}, id); result.Error != nil {
		ShowError(c, result.Error)
		return
	}

	// Redirect to settings.
	c.Redirect(http.StatusFound, "/settings")
}
