package controllers

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func ShowError(w http.ResponseWriter, statusCode int, err error) {
	log.WithField("code", statusCode).Error(err)

	w.WriteHeader(statusCode)
	Render("error.html.tmpl", w, map[string]any{
		"title":       "Error",
		"description": "The request failed.",
		"error":       err.Error(),
	})
}
