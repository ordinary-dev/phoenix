package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ordinary-dev/phoenix/backend"
	"github.com/ordinary-dev/phoenix/views"
	"log"
)

func main() {
	db, err := backend.GetDatabaseConnection()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./assets")

	// Main page
	r.GET("/", func(c *gin.Context) {
		views.ShowMainPage(c, db)
	})

	r.GET("/settings", func(c *gin.Context) {
		views.ShowSettings(c, db)
	})

	// Create new user
	r.POST("/users", func(c *gin.Context) {
		views.CreateUser(c, db)
	})

	r.POST("/signin", func(c *gin.Context) {
		views.AuthorizeUser(c, db)
	})

	// Create new group
	r.POST("/groups", func(c *gin.Context) {
		views.CreateGroup(c, db)
	})

	// Create new link
	r.POST("/links", func(c *gin.Context) {
		views.CreateLink(c, db)
	})

	// Update link.
	// HTML forms cannot be submitted using PUT or PATCH methods without javascript.
	r.POST("/links/:id/put", func(c *gin.Context) {
		views.UpdateLink(c, db)
	})

	// Delete link
	// HTML forms cannot be submitted using the DELETE method without javascript.
	r.POST("/links/:id/delete", func(c *gin.Context) {
		views.DeleteLink(c, db)
	})

	r.Run()
}
