package views

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ShowError(ctx *gin.Context, err error) {
	ctx.HTML(
		http.StatusBadRequest,
		"error.html.tmpl",
		gin.H{
			"error": err.Error(),
		},
	)
	ctx.Abort()
}
