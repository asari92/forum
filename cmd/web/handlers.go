package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"forum/internal/models"
	"forum/internal/validator"
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

func (app *application) postCreateView(w http.ResponseWriter, r *http.Request) {
	categories, err := app.categories.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Categories = categories
	data.Form = postCreateForm{
		// первая категория всегда отмечена
		Categories: []int{models.DefaultCategory},
	}

	app.render(w, http.StatusOK, "create_post.html", data)
}

type postCreateForm struct {
	Title      string
	Content    string
	Categories []int
	validator.Validator
}

func (app *application) postCreate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

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

	form := postCreateForm{
		Title:      r.PostForm.Get("title"),
		Content:    r.PostForm.Get("content"),
		Categories: categoryIDs,
	}

	// валидировать все данные
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")

	allCategories, err := app.categories.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	form.validateCategories(allCategories)

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Categories = allCategories
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create_post.html", data)
		return
	}

	//!!!!!!!!!!!!!!пока везде юзерайди = 1
	postId, err := app.posts.InsertPostWithCategories(form.Title, form.Content, 1, form.Categories)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/view/%d", postId), http.StatusSeeOther)
}

func (form *postCreateForm) validateCategories(allCategories []*models.Category) {
	var categoryIDs = []int{}
	if len(form.Categories) == 0 {
		form.AddFieldError("categories", "Need one or more category")
	} else {
		categoryMap := make(map[int]bool)
		for _, category := range allCategories {
			categoryMap[category.ID] = true
		}

		for _, c := range form.Categories {
			if _, exists := categoryMap[c]; exists {
				categoryIDs = append(categoryIDs, c)
			} else {
				form.AddFieldError("categories", "One or more categories are invalid")
			}
		}
	}

	if len(categoryIDs) > 0 {
		form.Categories = categoryIDs
	} else {
		// первая категория будет отмечена по умолчанию
		form.Categories = []int{models.DefaultCategory}
	}
}
