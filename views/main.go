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

	engine.GET("/signin", ShowLoginForm(cfg))
	engine.POST("/api/users/signin", AuthorizeUser(db, cfg))

	engine.GET("/registration", ShowRegistrationForm(cfg, db))
	engine.POST("/api/users", CreateUser(db, cfg))

	// This group requires authorization before viewing.
	protected := engine.Group("/")
	protected.Use(AuthMiddleware(db, cfg))

	// Main page
	protected.GET("/", ShowMainPage(cfg, db))

	protected.GET("/settings", ShowSettings(cfg, db))

	// Create new group
	protected.POST("/api/groups", CreateGroup(cfg, db))

	// Update group
	// HTML forms cannot be submitted using PUT or PATCH methods without javascript.
	protected.POST("/api/groups/:id/put", UpdateGroup(cfg, db))

	// Delete group
	// HTML forms cannot be submitted using the DELETE method without javascript.
	protected.POST("/api/groups/:id/delete", DeleteGroup(cfg, db))

	// Create new link
	protected.POST("/api/links", CreateLink(cfg, db))

	// Update link.
	// HTML forms cannot be submitted using PUT or PATCH methods without javascript.
	protected.POST("/api/links/:id/put", UpdateLink(cfg, db))

	// Delete link
	// HTML forms cannot be submitted using the DELETE method without javascript.
	protected.POST("/api/links/:id/delete", DeleteLink(cfg, db))

	return engine
}
