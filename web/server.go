package web

import (
	"net/http"

	"github.com/ordinary-dev/phoenix/web/controllers"
	"github.com/ordinary-dev/phoenix/web/middleware"
)

// Create and configure an HTTP server.
func GetHandler() (http.Handler, error) {
	if err := controllers.LoadTemplates(); err != nil {
		return nil, err
	}

	mux := http.NewServeMux()

	mux.Handle("/assets/", http.StripPrefix(
		"/assets",
		http.FileServer(http.Dir("web/assets")),
	))

	mux.HandleFunc("GET /signin", controllers.ShowSignInForm)
	mux.HandleFunc("POST /signin", controllers.AuthorizeUser)

	mux.HandleFunc("GET /registration", controllers.ShowRegistrationForm)
	mux.HandleFunc("POST /registration", controllers.CreateUser)

	protectedMux := http.NewServeMux()
	mux.Handle("/", middleware.RequireAuth(protectedMux))

	protectedMux.HandleFunc("GET /", controllers.ShowMainPage)
	protectedMux.HandleFunc("GET /settings", controllers.ShowSettings)

	// Groups.
	protectedMux.HandleFunc("POST /groups", controllers.CreateGroup)
	protectedMux.HandleFunc("POST /groups/{id}/update", controllers.UpdateGroup)
	protectedMux.HandleFunc("POST /groups/{id}/delete", controllers.DeleteGroup)

	// Links.
	protectedMux.HandleFunc("POST /links", controllers.CreateLink)
	protectedMux.HandleFunc("POST /links/{id}/update", controllers.UpdateLink)
	protectedMux.HandleFunc("POST /links/{id}/delete", controllers.DeleteLink)

	// Import-export
	protectedMux.HandleFunc("GET /export", controllers.Export)
	protectedMux.HandleFunc("GET /import", controllers.ImportPage)
	protectedMux.HandleFunc("POST /import", controllers.Import)

	return middleware.LoggingMiddleware(
		middleware.SecurityHeadersMiddleware(mux),
	), nil
}
