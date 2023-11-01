package views

import (
	"github.com/gin-gonic/gin"
	"github.com/ordinary-dev/phoenix/database"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func CreateGroup(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Save new group to the database.
		group := database.Group{
			Name: ctx.PostForm("groupName"),
		}
		if result := db.Create(&group); result.Error != nil {
			ShowError(ctx, result.Error)
			return
		}

		// This page is called from the settings, return the user back.
		ctx.Redirect(http.StatusFound, "/settings")
	}
}

func UpdateGroup(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
		if err != nil {
			ShowError(ctx, err)
			return
		}

		var group database.Group
		if result := db.First(&group, id); result.Error != nil {
			ShowError(ctx, result.Error)
			return
		}

		group.Name = ctx.PostForm("groupName")
		if result := db.Save(&group); result.Error != nil {
			ShowError(ctx, result.Error)
			return
		}

		// This page is called from the settings, return the user back.
		ctx.Redirect(http.StatusFound, "/settings")
	}
}

func DeleteGroup(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
		if err != nil {
			ShowError(ctx, err)
			return
		}

		if result := db.Delete(&database.Group{}, id); result.Error != nil {
			ShowError(ctx, result.Error)
			return
		}

		// Redirect to settings.
		ctx.Redirect(http.StatusFound, "/settings")
	}
}
