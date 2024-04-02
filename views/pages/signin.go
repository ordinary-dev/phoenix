package pages

import (
	"net/http"
	"strings"

	"github.com/ordinary-dev/phoenix/database"
	"github.com/ordinary-dev/phoenix/jwttoken"
)

func ShowSignInForm(w http.ResponseWriter, _ *http.Request) {
	Render("auth.html.tmpl", w, map[string]any{
		"title":       "Sign in",
		"description": "Authorization is required to view this page.",
		"button":      "Sign in",
		"formAction":  "/signin",
	})
}

func AuthorizeUser(w http.ResponseWriter, r *http.Request) {
	// Check credentials.
	username := strings.TrimSpace(r.FormValue("username"))
	password := strings.TrimSpace(r.FormValue("password"))
	_, err := database.GetAdminIfPasswordMatches(username, password)
	if err != nil {
		ShowError(w, http.StatusUnauthorized, err)
		return
	}

	// Generate an access token.
	token, err := jwttoken.GetJWTToken()
	if err != nil {
		ShowError(w, http.StatusInternalServerError, err)
		return
	}
	http.SetCookie(w, jwttoken.TokenToCookie(token))

	http.Redirect(w, r, "/", http.StatusFound)
}
