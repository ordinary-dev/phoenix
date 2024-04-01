package views

import (
	"net/http"

	"github.com/ordinary-dev/phoenix/views/middleware"
	"github.com/ordinary-dev/phoenix/views/pages"
)

// Create and configure an HTTP server.
//
// Unfortunately, I haven't found a way to use PUT and DELETE methods without JavaScript.
// POST is used instead.
func GetHttpServer() (*http.Server, error) {
	if err := pages.LoadTemplates(); err != nil {
		return nil, err
	}

	mux := http.NewServeMux()

	mux.Handle("/assets/", http.StripPrefix(
		"/assets",
		http.FileServer(http.Dir("assets")),
	))

	mux.HandleFunc("GET /signin", pages.ShowSignInForm)
	mux.HandleFunc("POST /signin", pages.AuthorizeUser)

	mux.HandleFunc("GET /registration", pages.ShowRegistrationForm)
	mux.HandleFunc("POST /registration", pages.CreateUser)

	protectedMux := http.NewServeMux()
	mux.Handle("/", middleware.RequireAuth(protectedMux))

	protectedMux.HandleFunc("GET /", pages.ShowMainPage)
	protectedMux.HandleFunc("GET /settings", pages.ShowSettings)

	// Groups.

	// Create group.
	protectedMux.HandleFunc("POST /groups", pages.CreateGroup)
	// Update group.
	protectedMux.HandleFunc("POST /groups/{id}/update", pages.UpdateGroup)
	// Delete group.
	protectedMux.HandleFunc("POST /groups/{id}/delete", pages.DeleteGroup)

	// Links.

	// Create link.
	protectedMux.HandleFunc("POST /links", pages.CreateLink)
	// Update link.
	protectedMux.HandleFunc("POST /links/{id}/update", pages.UpdateLink)
	// Delete link.
	protectedMux.HandleFunc("POST /links/{id}/delete", pages.DeleteLink)

	// Import-export
	protectedMux.HandleFunc("GET /export", pages.Export)
	protectedMux.HandleFunc("GET /import", pages.ImportPage)
	protectedMux.HandleFunc("POST /import", pages.Import)

	return &http.Server{
		Addr: ":8080",
		Handler: middleware.LoggingMiddleware(
			middleware.SecurityHeadersMiddleware(mux),
		),
	}, nil
}
