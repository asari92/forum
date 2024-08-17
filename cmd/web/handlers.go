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
	categories, err := app.categories.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Categories = categories

	app.render(w, http.StatusOK, "create_post.html", data)
}

func (app *application) postCreate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")

	// пока непонятно откуда брать юзерайди кажется из сессии
	// userID, err := strconv.Atoi(r.PostForm.Get("userID"))
	// if err != nil {
	// 	app.clientError(w, http.StatusBadRequest)
	// 	return
	// }

	var categoryIDs []int
	for _, id := range r.PostForm["categories"] {
		intID, err := strconv.Atoi(id)
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		categoryIDs = append(categoryIDs, intID)
	}

	// валидировать все данные

	//!!!!!!!!!!!!!!пока везде юзерайди = 1
	postId, err := app.posts.InsertPostWithCategories(title, content, 1, categoryIDs)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/view/%d", postId), http.StatusSeeOther)
}
