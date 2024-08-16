package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"forum/internal/models"
)

func (app *application) errorHandler(w http.ResponseWriter, r *http.Request) {
	app.notFound(w) // Use the notFound() helper
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	posts, err := app.posts.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Posts = posts

	app.render(w, http.StatusOK, "home.html", data)
}

func (app *application) postView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	post, err := app.posts.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Post = post

	app.render(w, http.StatusOK, "view.html", data)

	// возможный способ создания куки
	// session, _ := h.service.CreateSession(models)
	// cokkies := &http.Cookie{
	// 	Value: session.ID,
	// }

	// http.SetCookie(w, cokkies)
}

func (app *application) postCreateForm(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display the form for creating a new snippet..."))
}

func (app *application) postCreate(w http.ResponseWriter, r *http.Request) {
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	userID := 1

	id, err := app.posts.Insert(title, content, userID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
