package main

import (
	"html/template"
	"io/fs"
	"path/filepath"

	"forum/internal/models"
	"forum/ui"
)

type ReactionData struct {
	Likes        int
	Dislikes     int
	UserReaction *models.PostReaction
}

type templateData struct {
	CurrentYear     int
	Categories      []*models.Category
	CSRFToken       string
	Post            *models.Post
	Posts           []*models.Post
	Comments        []*models.Comment
	Form            any
	Flash           string
	IsAuthenticated bool
	User            *models.User
	ReactionData    *ReactionData
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
	"contains": contains,
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
