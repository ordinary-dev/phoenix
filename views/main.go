package views

import (
	"github.com/gin-gonic/gin"
	"github.com/ordinary-dev/phoenix/config"
	"gorm.io/gorm"
)

func GetGinEngine(cfg *config.Config, db *gorm.DB) *gin.Engine {
	engine := gin.Default()
	engine.LoadHTMLGlob("templates/*")
	engine.Static("/assets", "./assets")

	engine.Use(SecurityHeadersMiddleware)

	engine.GET("/signin", func(c *gin.Context) {
		ShowLoginForm(c)
	})
	engine.POST("/signin", func(c *gin.Context) {
		AuthorizeUser(c, db, cfg)
	})

	protected := engine.Group("/")
	protected.Use(func(c *gin.Context) {
		AuthMiddleware(c, cfg)
	})

	// Main page
	protected.GET("/", func(c *gin.Context) {
		ShowMainPage(c, db)
	})

	protected.GET("/settings", func(c *gin.Context) {
		ShowSettings(c, db)
	})

	// Create new user
	protected.POST("/users", func(c *gin.Context) {
		CreateUser(c, db, cfg)
	})

	// Create new group
	protected.POST("/groups", func(c *gin.Context) {
		CreateGroup(c, db)
	})

	// Update group
	// HTML forms cannot be submitted using PUT or PATCH methods without javascript.
	protected.POST("/groups/:id/put", func(c *gin.Context) {
		UpdateGroup(c, db)
	})

	// Delete group
	// HTML forms cannot be submitted using the DELETE method without javascript.
	protected.POST("/groups/:id/delete", func(c *gin.Context) {
		DeleteGroup(c, db)
	})

	// Create new link
	protected.POST("/links", func(c *gin.Context) {
		CreateLink(c, db)
	})

	// Update link.
	// HTML forms cannot be submitted using PUT or PATCH methods without javascript.
	protected.POST("/links/:id/put", func(c *gin.Context) {
		UpdateLink(c, db)
	})

	// Delete link
	// HTML forms cannot be submitted using the DELETE method without javascript.
	protected.POST("/links/:id/delete", func(c *gin.Context) {
		DeleteLink(c, db)
	})

	return engine
}
