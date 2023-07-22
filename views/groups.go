package views

import (
	"github.com/gin-gonic/gin"
	"github.com/ordinary-dev/phoenix/database"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func CreateGroup(c *gin.Context, db *gorm.DB) {
	// Save new group to the database.
	group := database.Group{
		Name: c.PostForm("groupName"),
	}
	if result := db.Create(&group); result.Error != nil {
		ShowError(c, result.Error)
		return
	}

	// This page is called from the settings, return the user back.
	c.Redirect(http.StatusFound, "/settings")
}

func UpdateGroup(c *gin.Context, db *gorm.DB) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		ShowError(c, err)
		return
	}

	var group database.Group
	if result := db.First(&group, id); result.Error != nil {
		ShowError(c, result.Error)
		return
	}

	group.Name = c.PostForm("groupName")
	if result := db.Save(&group); result.Error != nil {
		ShowError(c, result.Error)
		return
	}

	// This page is called from the settings, return the user back.
	c.Redirect(http.StatusFound, "/settings")
}

func DeleteGroup(c *gin.Context, db *gorm.DB) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		ShowError(c, err)
		return
	}

	if result := db.Delete(&database.Group{}, id); result.Error != nil {
		ShowError(c, result.Error)
		return
	}

	// Redirect to settings.
	c.Redirect(http.StatusFound, "/settings")
}
