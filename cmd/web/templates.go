package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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
	Form            any
	Flash           string
	IsAuthenticated bool
	User            *models.User
	ReactionData    *ReactionData
	CurrentPage     int
	PageSize        int
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

func join(items []int, sep rune) string {
	strItems := make([]string, len(items))
	for i, item := range items {
		strItems[i] = strconv.Itoa(item)
	}
	return strings.Join(strItems, string(sep))
}

var functions = template.FuncMap{
	"humanDate": humanDate,
	"contains":  contains,
	"add":       func(a, b int) int { return a + b },
	"sub":       func(a, b int) int { return a - b },
	"join":      join,
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
