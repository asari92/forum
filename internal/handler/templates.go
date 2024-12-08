package handler

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"

	"forum/internal/entities"
	"forum/ui"
)

type AppError struct {
	Message    string
	StatusCode int
}

type ReactionData struct {
	Likes        int
	Dislikes     int
	UserReaction *entities.PostReaction
}

func (rd *ReactionData) GetUserReaction() int {
	if rd.UserReaction != nil {
		if rd.UserReaction.IsLike {
			return 1
		} else {
			return -1
		}
	}
	return 0
}

type templateData struct {
	AppError        AppError
	CurrentYear     int
	Categories      []*entities.Category
	CSRFToken       string
	Post            *entities.Post
	Posts           []*entities.Post
	Images          []*entities.Image
	ImageError      bool
	Comments        []*entities.Comment
	Form            any
	Flash           string
	IsAuthenticated bool
	User            *entities.User
	ReactionData    *ReactionData
	Header          string
	Pagination      any
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
	"contains": contains,
	"add":      func(a, b int) int { return a + b },
	"sub":      func(a, b int) int { return a - b },
	"join":     join,
}

func NewTemplateCache() (map[string]*template.Template, error) {
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
