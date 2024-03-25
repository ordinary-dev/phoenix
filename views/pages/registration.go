package pages

import (
	"errors"
	"net/http"
	"strings"

	"github.com/ordinary-dev/phoenix/database"
	"github.com/ordinary-dev/phoenix/jwttoken"
)

func ShowRegistrationForm(w http.ResponseWriter, _ *http.Request) {
	userCount, err := database.CountAdmins()
	if err != nil {
		ShowError(w, http.StatusInternalServerError, err)
		return
	}

	if userCount > 0 {
		ShowError(w, http.StatusBadRequest, errors.New("at least 1 user already exists"))
		return
	}

	Render("auth.html.tmpl", w, map[string]any{
		"title":       "Create an account",
		"description": "To prevent other people from seeing your links, create an account.",
		"button":      "Create",
		"formAction":  "/registration",
	})
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	userCount, err := database.CountAdmins()
	if err != nil {
		ShowError(w, http.StatusInternalServerError, err)
		return
	}

	if userCount > 0 {
		ShowError(w, http.StatusBadRequest, errors.New("at least 1 user already exists"))
		return
	}

	// Try to create a user.
	username := strings.TrimSpace(r.FormValue("username"))
	password := strings.TrimSpace(r.FormValue("password"))
	_, err = database.CreateAdmin(username, password)
	if err != nil {
		ShowError(w, http.StatusInternalServerError, err)
		return
	}

	// Generate access token.
	token, err := jwttoken.GetJWTToken()
	if err != nil {
		ShowError(w, http.StatusInternalServerError, err)
		return
	}
	http.SetCookie(w, jwttoken.TokenToCookie(token))

	// Redirect to homepage.
	http.Redirect(w, r, "/", http.StatusFound)
}
