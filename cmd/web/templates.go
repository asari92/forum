package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"forum/internal/models"
	"forum/ui"
)

type templateData struct {
	CurrentYear     int
	Categories      []*models.Category
	CSRFToken       string
	Post            *models.Post
	Posts           []*models.Post
	Form            any
	Flash           string
	IsAuthenticated bool
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

var functions = template.FuncMap{
	"humanDate": humanDate,
	"contains":  contains,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.html",
			"html/partials/*.html",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
