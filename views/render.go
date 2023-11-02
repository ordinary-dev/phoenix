package views

import (
	"github.com/gin-gonic/gin"
	"github.com/ordinary-dev/phoenix/config"
)

// Fill in the necessary parameters from the settings and output html.
func Render(ctx *gin.Context, cfg *config.Config, status int, templatePath string, params map[string]any) {
	params["WebsiteTitle"] = cfg.Title
	params["FontFamily"] = cfg.FontFamily
	ctx.HTML(status, templatePath, params)
}
