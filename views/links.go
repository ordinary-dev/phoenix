package views

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ordinary-dev/phoenix/config"
	"github.com/ordinary-dev/phoenix/database"
	"gorm.io/gorm"
)

func CreateLink(cfg *config.Config, db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		groupID, err := strconv.ParseUint(ctx.PostForm("groupID"), 10, 32)
		if err != nil {
			ShowError(ctx, cfg, err)
			return
		}

		link := database.Link{
			Name:    ctx.PostForm("linkName"),
			Href:    ctx.PostForm("href"),
			GroupID: groupID,
		}
		icon := ctx.PostForm("icon")
		if icon == "" {
			link.Icon = nil
		} else {
			link.Icon = &icon
		}
		if result := db.Create(&link); result.Error != nil {
			ShowError(ctx, cfg, result.Error)
			return
		}

		// Redirect to settings.
		ctx.Redirect(http.StatusFound, fmt.Sprintf("/settings#link-%v", link.ID))
	}
}

func UpdateLink(cfg *config.Config, db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
		if err != nil {
			ShowError(ctx, cfg, err)
			return
		}

		var link database.Link
		if result := db.First(&link, id); result.Error != nil {
			ShowError(ctx, cfg, err)
			return
		}

		link.Name = ctx.PostForm("linkName")
		link.Href = ctx.PostForm("href")
		icon := ctx.PostForm("icon")
		if icon == "" {
			link.Icon = nil
		} else {
			link.Icon = &icon
		}
		if result := db.Save(&link); result.Error != nil {
			ShowError(ctx, cfg, result.Error)
			return
		}

		// Redirect to settings.
		ctx.Redirect(http.StatusFound, fmt.Sprintf("/settings#link-%v", link.ID))
	}
}

func DeleteLink(cfg *config.Config, db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
		if err != nil {
			ShowError(ctx, cfg, err)
			return
		}

		if result := db.Delete(&database.Link{}, id); result.Error != nil {
			ShowError(ctx, cfg, result.Error)
			return
		}

		// Redirect to settings.
		ctx.Redirect(http.StatusFound, "/settings")
	}
}
