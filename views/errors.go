package views

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ShowError(c *gin.Context, err error) {
	c.HTML(
		http.StatusBadRequest,
		"error.html.tmpl",
		gin.H{
			"error": err.Error(),
		},
	)
}
