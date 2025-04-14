package web

import (
	"net/http"

	"github.com/ordinary-dev/phoenix/database"
	"github.com/ordinary-dev/phoenix/web/controllers"
	"github.com/ordinary-dev/phoenix/web/middleware"
)

// Create and configure an HTTP server.
func GetHttpHandler(db database.Database) (http.Handler, error) {
	ctrl, err := controllers.New(db)
	if err != nil {
		return nil, err
	}

	middlewareInstance := middleware.New(ctrl, db)

	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix(
		"/assets",
		http.FileServer(http.Dir("web/assets")),
	))

	mux.HandleFunc("GET /signin", ctrl.ShowSignInForm)
	mux.HandleFunc("POST /signin", ctrl.AuthorizeUser)

	mux.HandleFunc("GET /registration", ctrl.ShowRegistrationForm)
	mux.HandleFunc("POST /registration", ctrl.CreateUser)

	protectedMux := http.NewServeMux()
	mux.Handle("/", middlewareInstance.RequireAuth(protectedMux))

	protectedMux.HandleFunc("GET /", ctrl.ShowMainPage)
	protectedMux.HandleFunc("GET /settings", ctrl.ShowSettings)

	// Groups.
	protectedMux.HandleFunc("POST /groups", ctrl.CreateGroup)
	protectedMux.HandleFunc("POST /groups/{id}/update", ctrl.UpdateGroup)
	protectedMux.HandleFunc("POST /groups/{id}/delete", ctrl.DeleteGroup)

	// Links.
	protectedMux.HandleFunc("POST /links", ctrl.CreateLink)
	protectedMux.HandleFunc("POST /links/{id}/update", ctrl.UpdateLink)
	protectedMux.HandleFunc("POST /links/{id}/delete", ctrl.DeleteLink)

	// Import-export
	protectedMux.HandleFunc("GET /export", ctrl.Export)
	protectedMux.HandleFunc("GET /import", ctrl.ImportPage)
	protectedMux.HandleFunc("POST /import", ctrl.Import)

	return middlewareInstance.LoggingMiddleware(
		middlewareInstance.SecurityHeadersMiddleware(mux),
	), nil
}
