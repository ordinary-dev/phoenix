package controllers

import (
	"html/template"
	"io"
	"log/slog"
	"os"
	"path"

	"github.com/ordinary-dev/phoenix/config"
	"github.com/ordinary-dev/phoenix/database"
)

type Controllers struct {
	db database.Database
	// Preloaded templates.
	// The key is the file name.
	templates map[string]*template.Template
}

func New(db database.Database) (*Controllers, error) {
	c := Controllers{
		db: db,
	}

	var err error
	c.templates, err = loadTemplates()
	return &c, err
}

// Preload all templates into `Templates` map.
// Map key is the file name.
func loadTemplates() (map[string]*template.Template, error) {
	templates := make(map[string]*template.Template)

	// Fragments are reusable parts of templates.
	fragments, err := os.ReadDir("web/views/fragments")
	if err != nil {
		return nil, err
	}

	var fragmentPaths []string
	for _, f := range fragments {
		fragmentPaths = append(
			fragmentPaths,
			path.Join("web/views/fragments", f.Name()),
		)
	}

	files, err := os.ReadDir("web/views")
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		templatePaths := []string{path.Join("web/views", f.Name())}
		templatePaths = append(templatePaths, fragmentPaths...)

		tmpl, err := template.ParseFiles(templatePaths...)
		if err != nil {
			return nil, err
		}

		templates[f.Name()] = tmpl
		slog.Debug("template was loaded", "file", f.Name())
	}

	return templates, nil
}

func (c *Controllers) render(template string, wr io.Writer, data map[string]any) {
	data["fontFamily"] = config.Cfg.FontFamily

	if _, ok := data["title"]; !ok {
		data["title"] = config.Cfg.Title
	}

	err := c.templates[template].Execute(wr, data)
	if err != nil {
		slog.Error("template rendering failed", "err", err, "template", template)
	}
}
