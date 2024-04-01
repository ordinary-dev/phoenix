package pages

import (
	"net/http"
	"time"

	"github.com/ordinary-dev/phoenix/database"
	log "github.com/sirupsen/logrus"
)

func ShowMainPage(w http.ResponseWriter, r *http.Request) {
	groups, err := database.GetGroupsWithLinks()
	if err != nil {
		ShowError(w, http.StatusInternalServerError, err)
		return
	}

	// Get desired style.
	style := r.FormValue("style")

	if style == "tiles" || style == "list" {
		// If a valid style is specified in the url -
		// save the value in a cookie.
		http.SetCookie(w, &http.Cookie{
			Name:    "phoenix-style",
			Value:   style,
			Expires: time.Now().Add(time.Hour * 24 * 30 * 12 * 10),
		})
	} else {
		// The style is not specified or has an incorrect type, check the cookies.
		styleCookie, err := r.Cookie("phoenix-style")
		if err == nil {
			style = styleCookie.Value
		}
	}

	if style != "tiles" && style != "list" {
		style = "list"
	}

	err = Render("index.html.tmpl", w, map[string]any{
		"description": "Self-hosted start page.",
		"groups":      groups,
		"style":       style,
	})
	if err != nil {
		log.Error(err)
	}
}
