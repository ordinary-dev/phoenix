package views

import (
	"github.com/gin-gonic/gin"
	"github.com/ordinary-dev/phoenix/database"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func CreateLink(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		groupID, err := strconv.ParseUint(ctx.PostForm("groupID"), 10, 32)
		if err != nil {
			ShowError(ctx, err)
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
			ShowError(ctx, result.Error)
			return
		}

		// Redirect to settings.
		ctx.Redirect(http.StatusFound, "/settings")
	}
}

func UpdateLink(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
		if err != nil {
			ShowError(ctx, err)
			return
		}

		var link database.Link
		if result := db.First(&link, id); result.Error != nil {
			ShowError(ctx, err)
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
			ShowError(ctx, result.Error)
			return
		}

		// Redirect to settings.
		ctx.Redirect(http.StatusFound, "/settings")
	}
}

func DeleteLink(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
		if err != nil {
			ShowError(ctx, err)
			return
		}

		if result := db.Delete(&database.Link{}, id); result.Error != nil {
			ShowError(ctx, result.Error)
			return
		}

		// Redirect to settings.
		ctx.Redirect(http.StatusFound, "/settings")
	}
}
