package controllers

import (
	"log/slog"
	"net/http"
)

func ShowError(w http.ResponseWriter, statusCode int, err error) {
	if statusCode < 500 {
		slog.Warn("request failed", "err", err, "code", statusCode)
	} else {
		slog.Error("request failed", "err", err, "code", statusCode)
	}

	w.WriteHeader(statusCode)
	Render("error.html.tmpl", w, map[string]any{
		"title":       "Error",
		"description": "The request failed.",
		"error":       err.Error(),
	})
}
