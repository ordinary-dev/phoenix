package controllers

import (
	"html/template"
	"io"
	"os"
	"path"

	log "github.com/sirupsen/logrus"

	"github.com/ordinary-dev/phoenix/config"
)

var (
	// Preloaded templates.
	// The key is the file name.
	templates = make(map[string]*template.Template)
)

// Preload all templates into `Templates` map.
func LoadTemplates() error {
	// Fragments are reusable parts of templates.
	fragments, err := os.ReadDir("web/views/fragments")
	if err != nil {
		return err
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
		return err
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		templatePaths := []string{path.Join("web/views", f.Name())}
		templatePaths = append(templatePaths, fragmentPaths...)

		tmpl, err := template.ParseFiles(templatePaths...)
		if err != nil {
			return err
		}

		templates[f.Name()] = tmpl
		log.Infof("Template %v was loaded", f.Name())
	}

	return nil
}

func Render(template string, wr io.Writer, data map[string]any) {
	data["fontFamily"] = config.Cfg.FontFamily

	if _, ok := data["title"]; !ok {
		data["title"] = config.Cfg.Title
	}

	err := templates[template].Execute(wr, data)
	if err != nil {
		log.WithField("template", template).Error(err)
	}
}
