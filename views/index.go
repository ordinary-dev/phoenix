package views

import (
	"github.com/gin-gonic/gin"
	"github.com/ordinary-dev/phoenix/backend"
	"gorm.io/gorm"
	"net/http"
)

func ShowMainPage(c *gin.Context, db *gorm.DB) {
	if err := RequireAuth(c, db); err != nil {
		return
	}

	// Get a list of groups with links
	groups, err := backend.GetGroups(db)
	if err != nil {
		ShowError(c, err)
		return
	}

	c.HTML(http.StatusOK, "index.html.tmpl", gin.H{
		"groups": groups,
	})
}
