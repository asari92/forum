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

	data := app.newTemplateData(w, r)
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

	data := app.newTemplateData(w, r)
	data.Post = post

	app.render(w, http.StatusOK, "post_view.html", data)

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

	data := app.newTemplateData(w, r)
	// Извлечение CSRF-токена из контекста
	token := r.Context().Value("csrfToken").(string)
	data.CSRFToken = token
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
		data := app.newTemplateData(w, r)
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

	err = app.sessionManager.SessionStart(w, r).Set("flash", "Post successfully created!")
	if err != nil {
		app.serverError(w, err)
	}

	http.Redirect(w, r, fmt.Sprintf("/post/view/%d", postId), http.StatusSeeOther)
}

func (form *postCreateForm) validateCategories(allCategories []*models.Category) {
	categoryIDs := []int{}
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

func (app *application) userSignupView(w http.ResponseWriter, r *http.Request) {
	// sess := app.sessionManager.SessionStart(w, r)
	data := app.newTemplateData(w, r)
	// Извлечение CSRF-токена из контекста
	token := r.Context().Value("csrfToken").(string)
	data.CSRFToken = token
	// data.Session = sess.Get("username")                    !!!!!!!!!!!!!!!!!!! ДЛЯ ЧЕГО Я ЭТО ДЕЛАЛ?!
	data.Form = signupForm{}
	app.render(w, http.StatusOK, "signup.html", data)
}

type signupForm struct {
	Username string
	Email    string
	Password string
	validator.Validator
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form := signupForm{
		Username: r.PostForm.Get("username"),
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}

	// ДОБАВИТЬ ПРОВЕРОК!!!!!!!!!!!!!!!!!!!
	form.CheckField(validator.NotBlank(form.Username), "username", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Username, 100), "username", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.MaxChars(form.Email, 100), "email", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")
	form.CheckField(validator.MaxChars(form.Password, 100), "password", "This field cannot be more than 100 characters long")

	if !form.Valid() {
		data := app.newTemplateData(w, r)
		data.Form = form
		data.CSRFToken = r.Context().Value("csrfToken").(string)
		app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
		return
	}

	err = app.users.Insert(form.Username, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")

			data := app.newTemplateData(w, r)
			data.Form = form
			data.CSRFToken = r.Context().Value("csrfToken").(string)
			app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
		} else {
			app.serverError(w, err)
		}

		return
	}
	sess := app.sessionManager.SessionStart(w, r)
	// Если валидация прошла успешно, удаляем токен из сессии
	sess.Delete("token")

	err = sess.Set("flash", "Your signup was successful. Please log in.")
	if err != nil {
		app.serverError(w, err)
	}

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userLoginView(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Display a HTML form for logging in a user...")
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Authenticate and login the user...")

	// // Проверка имени пользователя и пароля
	// // if !app.authenticate(username, password) {
	// // 	// Если проверка не удалась, вернуть сообщение об ошибке
	// // 	http.Error(w, "Invalid credentials", http.StatusUnauthorized)
	// // 	return
	// // }

	// // Получение роли пользователя из базы данных или другого источника
	// // role := app.users.getUserRole(username)

	// // Запуск сессии и сохранение данных пользователя
	// sess := app.sessionManager.SessionStart(w, r)
	// sess.Set("username", form.Username)
	// // Если валидация прошла успешно, удаляем токен из сессии
	// sess.Delete("token")
	// // sess.Set("role", role) // Сохраняем роль в сессии
	// // Перенаправление на главную страницу после успешного входа
	// http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) userLogout(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout the user...")
}
