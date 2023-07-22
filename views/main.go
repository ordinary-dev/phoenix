package views

import (
	"github.com/gin-gonic/gin"
	"github.com/ordinary-dev/phoenix/config"
	"gorm.io/gorm"
)

func GetGinEngine(cfg *config.Config, db *gorm.DB) *gin.Engine {
	if cfg.Production {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(gin.Recovery())
	if cfg.EnableGinLogger {
		engine.Use(gin.Logger())
	}

	engine.LoadHTMLGlob("templates/*")
	engine.Static("/assets", "./assets")

	engine.Use(SecurityHeadersMiddleware)

	engine.GET("/signin", func(c *gin.Context) {
		ShowLoginForm(c)
	})
	engine.POST("/api/users/signin", func(c *gin.Context) {
		AuthorizeUser(c, db, cfg)
	})

	engine.GET("/registration", func(c *gin.Context) {
		ShowRegistrationForm(c, db)
	})
	engine.POST("/api/users", func(c *gin.Context) {
		CreateUser(c, db, cfg)
	})

	// This group requires authorization before viewing.
	protected := engine.Group("/")
	protected.Use(func(c *gin.Context) {
		AuthMiddleware(c, db, cfg)
	})

	// Main page
	protected.GET("/", func(c *gin.Context) {
		ShowMainPage(c, db)
	})

	protected.GET("/settings", func(c *gin.Context) {
		ShowSettings(c, db)
	})

	// Create new group
	protected.POST("/api/groups", func(c *gin.Context) {
		CreateGroup(c, db)
	})

	// Update group
	// HTML forms cannot be submitted using PUT or PATCH methods without javascript.
	protected.POST("/api/groups/:id/put", func(c *gin.Context) {
		UpdateGroup(c, db)
	})

	// Delete group
	// HTML forms cannot be submitted using the DELETE method without javascript.
	protected.POST("/api/groups/:id/delete", func(c *gin.Context) {
		DeleteGroup(c, db)
	})

	// Create new link
	protected.POST("/api/links", func(c *gin.Context) {
		CreateLink(c, db)
	})

	// Update link.
	// HTML forms cannot be submitted using PUT or PATCH methods without javascript.
	protected.POST("/api/links/:id/put", func(c *gin.Context) {
		UpdateLink(c, db)
	})

	// Delete link
	// HTML forms cannot be submitted using the DELETE method without javascript.
	protected.POST("/api/links/:id/delete", func(c *gin.Context) {
		DeleteLink(c, db)
	})

	return engine
}
