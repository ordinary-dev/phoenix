package views

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ordinary-dev/phoenix/config"
	"github.com/ordinary-dev/phoenix/database"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func CreateGroup(cfg *config.Config, db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Save new group to the database.
		group := database.Group{
			Name: ctx.PostForm("groupName"),
		}
		if result := db.Create(&group); result.Error != nil {
			ShowError(ctx, cfg, result.Error)
			return
		}

		// This page is called from the settings, return the user back.
		ctx.Redirect(http.StatusFound, fmt.Sprintf("/settings#group-%v", group.ID))
	}
}

func UpdateGroup(cfg *config.Config, db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
		if err != nil {
			ShowError(ctx, cfg, err)
			return
		}

		var group database.Group
		if result := db.First(&group, id); result.Error != nil {
			ShowError(ctx, cfg, result.Error)
			return
		}

		group.Name = ctx.PostForm("groupName")
		if result := db.Save(&group); result.Error != nil {
			ShowError(ctx, cfg, result.Error)
			return
		}

		// This page is called from the settings, return the user back.
		ctx.Redirect(http.StatusFound, fmt.Sprintf("/settings#group-%v", group.ID))
	}
}

func DeleteGroup(cfg *config.Config, db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
		if err != nil {
			ShowError(ctx, cfg, err)
			return
		}

		if result := db.Delete(&database.Group{}, id); result.Error != nil {
			ShowError(ctx, cfg, result.Error)
			return
		}

		// Redirect to settings.
		ctx.Redirect(http.StatusFound, "/settings")
	}
}
