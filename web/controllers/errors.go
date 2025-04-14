package controllers

import (
	"log/slog"
	"net/http"
)

func (c *Controllers) ShowError(w http.ResponseWriter, statusCode int, err error) {
	if statusCode < 500 {
		slog.Warn("request failed", "err", err, "code", statusCode)
	} else {
		slog.Error("request failed", "err", err, "code", statusCode)
	}

	w.WriteHeader(statusCode)
	c.render("error.html.tmpl", w, map[string]any{
		"title":       "Error",
		"description": "The request failed.",
		"error":       err.Error(),
	})
}
