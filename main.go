package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ordinary-dev/phoenix/backend"
	"github.com/ordinary-dev/phoenix/config"
	"github.com/ordinary-dev/phoenix/views"
	"github.com/sirupsen/logrus"
)

func main() {
	// Configure logger
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Read config
	cfg, err := config.GetConfig()
	if err != nil {
		logrus.Fatalf("%v", err)
	}

	// Set log level
	logLevel := cfg.GetLogLevel()
	logrus.SetLevel(logLevel)
	logrus.Infof("Setting log level to %v", logLevel)

	// Connect to the database
	db, err := backend.GetDatabaseConnection(cfg)
	if err != nil {
		logrus.Fatalf("%v", err)
	}

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./assets")

	r.GET("/signin", func(c *gin.Context) {
		views.ShowLoginForm(c)
	})
	r.POST("/signin", func(c *gin.Context) {
		views.AuthorizeUser(c, db, cfg)
	})

	protected := r.Group("/")
	protected.Use(func(c *gin.Context) {
		views.AuthMiddleware(c, cfg)
	})

	// Main page
	protected.GET("/", func(c *gin.Context) {
		views.ShowMainPage(c, db)
	})

	protected.GET("/settings", func(c *gin.Context) {
		views.ShowSettings(c, db)
	})

	// Create new user
	protected.POST("/users", func(c *gin.Context) {
		views.CreateUser(c, db, cfg)
	})

	// Create new group
	protected.POST("/groups", func(c *gin.Context) {
		views.CreateGroup(c, db)
	})

	// Update group
	// HTML forms cannot be submitted using PUT or PATCH methods without javascript.
	protected.POST("/groups/:id/put", func(c *gin.Context) {
		views.UpdateGroup(c, db)
	})

	// Delete group
	// HTML forms cannot be submitted using the DELETE method without javascript.
	protected.POST("/groups/:id/delete", func(c *gin.Context) {
		views.DeleteGroup(c, db)
	})

	// Create new link
	protected.POST("/links", func(c *gin.Context) {
		views.CreateLink(c, db)
	})

	// Update link.
	// HTML forms cannot be submitted using PUT or PATCH methods without javascript.
	protected.POST("/links/:id/put", func(c *gin.Context) {
		views.UpdateLink(c, db)
	})

	// Delete link
	// HTML forms cannot be submitted using the DELETE method without javascript.
	protected.POST("/links/:id/delete", func(c *gin.Context) {
		views.DeleteLink(c, db)
	})

	r.Run(":8080")
}
