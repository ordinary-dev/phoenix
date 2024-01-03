package views

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ordinary-dev/phoenix/config"
	"github.com/ordinary-dev/phoenix/database"
	"gorm.io/gorm"
)

func ShowSettings(cfg *config.Config, db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get a list of groups with links
		var groups []database.Group
		result := db.
			Model(&database.Group{}).
			Preload("Links").
			Find(&groups)

		if result.Error != nil {
			ShowError(ctx, cfg, result.Error)
			return
		}

		Render(ctx, cfg, http.StatusOK, "settings.html.tmpl", gin.H{
			"groups": groups,
		})
	}
}
