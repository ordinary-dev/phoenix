package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ordinary-dev/phoenix/backend"
	"github.com/ordinary-dev/phoenix/views"
	"log"
	"net/http"
	"strconv"
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
		if err := views.RequireAuth(c, db); err != nil {
			return
		}

		groups, err := backend.GetGroups(db)
		if err != nil {
			views.ShowError(c, err)
			return
		}
		c.HTML(http.StatusOK, "index.html.tmpl", gin.H{
			"groups": groups,
		})
	})

	// Settings
	r.GET("/settings", func(c *gin.Context) {
		if err := views.RequireAuth(c, db); err != nil {
			return
		}

		groups, err := backend.GetGroups(db)
		if err != nil {
			views.ShowError(c, err)
			return
		}

		c.HTML(http.StatusOK, "settings.html.tmpl", gin.H{
			"groups": groups,
		})
	})

	// Create new user
	r.POST("/users", func(c *gin.Context) {
		// If at least 1 administator exists, require authorization
		if backend.CountAdmins(db) > 0 {
			tokenValue, err := c.Cookie("phoenix-token")

			// Anonymous visitor
			if err != nil {
				err = errors.New("At least 1 user exists, you have to sign in first")
				views.ShowError(c, err)
				return
			}

			err = backend.ValidateToken(db, tokenValue)
			if err != nil {
				views.ShowError(c, err)
				return
			}
		}

		// User is authorized or no user exists.
		// Try to create a user.
		username := c.PostForm("username")
		password := c.PostForm("password")
		admin, err := backend.CreateAdmin(db, username, password)
		if err != nil {
			views.ShowError(c, err)
			return
		}

		// Generate access token.
		token, err := backend.CreateAccessToken(db, admin.ID)
		if err != nil {
			views.ShowError(c, err)
			return
		}
		backend.SetTokenCookie(c, token)

		// Redirect to homepage.
		c.Redirect(http.StatusFound, "/")
	})

	// Create new group
	r.POST("/groups", func(c *gin.Context) {
		if err := views.RequireAuth(c, db); err != nil {
			return
		}

		groupName := c.PostForm("groupName")
		_, err := backend.CreateGroup(db, groupName)
		if err != nil {
			views.ShowError(c, err)
			return
		}

		// Redirect to settings.
		c.Redirect(http.StatusFound, "/settings")
	})

	// Create new link
	r.POST("/links", func(c *gin.Context) {
		if err := views.RequireAuth(c, db); err != nil {
			return
		}

		linkName := c.PostForm("linkName")
		href := c.PostForm("href")
		groupID, err := strconv.ParseUint(c.PostForm("groupID"), 10, 32)
		if err != nil {
			views.ShowError(c, err)
			return
		}

		_, err = backend.CreateLink(db, linkName, href, groupID)
		if err != nil {
			views.ShowError(c, err)
			return
		}

		// Redirect to settings.
		c.Redirect(http.StatusFound, "/settings")
	})

	// Update link
	r.POST("/links/:id/put", func(c *gin.Context) {
		if err := views.RequireAuth(c, db); err != nil {
			return
		}

		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			views.ShowError(c, err)
			return
		}
		linkName := c.PostForm("linkName")
		href := c.PostForm("href")

		_, err = backend.UpdateLink(db, id, linkName, href)
		if err != nil {
			views.ShowError(c, err)
			return
		}

		// Redirect to settings.
		c.Redirect(http.StatusFound, "/settings")
	})

	// Delete link
	r.POST("/links/:id/delete", func(c *gin.Context) {
		if err := views.RequireAuth(c, db); err != nil {
			return
		}

		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			views.ShowError(c, err)
			return
		}

		err = backend.DeleteLink(db, id)
		if err != nil {
			views.ShowError(c, err)
			return
		}

		// Redirect to settings.
		c.Redirect(http.StatusFound, "/settings")
	})

	r.POST("/signin", func(c *gin.Context) {
		// Check credentials.
		username := c.PostForm("username")
		password := c.PostForm("password")
		admin, err := backend.AuthorizeAdmin(db, username, password)
		if err != nil {
			views.ShowError(c, err)
			return
		}

		// Generate an access token.
		token, err := backend.CreateAccessToken(db, admin.ID)
		if err != nil {
			views.ShowError(c, err)
			return
		}
		backend.SetTokenCookie(c, token)

		// Redirect to homepage.
		c.Redirect(http.StatusFound, "/")
	})

	r.Run()
}
