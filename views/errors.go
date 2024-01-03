package views

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ordinary-dev/phoenix/config"
)

func ShowError(ctx *gin.Context, cfg *config.Config, err error) {
	Render(ctx, cfg, http.StatusBadRequest, "error.html.tmpl", gin.H{
		"error": err.Error(),
	})
	ctx.Abort()
}
