package handler

import (
	"errors"
	"net/http"
	"strconv"

	"forum/internal/entities"
)

type pagination struct {
	CurrentPage      int
	HasNextPage      bool
	PaginationAction string
}

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	page := 1      // Определяем текущую страницу. По умолчанию - страница 1.
	pageSize := 10 // Количество постов на одной странице

	// Получаем посты для нужной страницы через юзкейс.
	userPostsDTO, err := app.Service.Post.GetAllPaginatedPostsDTO(page, pageSize, "/")
	if err != nil {
		app.Logger.Error("get all paginated posts", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	form := app.Service.Post.NewPostCreateForm()

	data := app.newTemplateData(r)
	data.Form = form
	data.Header = "All posts"
	data.Posts = userPostsDTO.Posts
	data.Categories = userPostsDTO.Categories
	data.Pagination = pagination{
		CurrentPage:      userPostsDTO.CurrentPage,
		HasNextPage:      userPostsDTO.HasNextPage,
		PaginationAction: userPostsDTO.PaginationURL, // Динамическая ссылка на пагинацию
	}
	app.render(w, http.StatusOK, "home.html", data)
}

func (app *Application) filterPosts(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.render(w, http.StatusNotFound, Errorpage, nil)
		return
	}
	err := r.ParseForm()
	if err != nil {
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return
	}

	page := 1
	pageSize := 10

	if p, err := strconv.Atoi(r.PostFormValue("page")); err == nil && p > 0 {
		page = p
	}

	// Получаем массив категорий
	categoryIDs := []int{}
	for _, id := range r.PostForm["categories"] {
		intID, err := strconv.Atoi(id)
		if err != nil {
			app.render(w, http.StatusBadRequest, Errorpage, nil)
			return
		}
		categoryIDs = append(categoryIDs, intID)
	}

	form := app.Service.Post.NewPostCreateForm()
	form.Categories = categoryIDs

	data := app.newTemplateData(r)
	data.Form = form

	filteredPostsDTO, err := app.Service.Post.GetFilteredPaginatedPostsDTO(&form, page, pageSize, "/")
	if filteredPostsDTO != nil {
		data.Posts = filteredPostsDTO.Posts
		data.Header = filteredPostsDTO.Header
		data.Categories = filteredPostsDTO.Categories
		data.Pagination = pagination{
			CurrentPage:      filteredPostsDTO.CurrentPage,
			HasNextPage:      filteredPostsDTO.HasNextPage,
			PaginationAction: filteredPostsDTO.PaginationURL, // Маршрут для пагинации
		}
	}
	if err != nil {
		app.Logger.Error("filter posts by categories", "error", err)
		if errors.Is(err, entities.ErrInvalidCredentials) {
			app.render(w, http.StatusUnprocessableEntity, "home.html", data)
		} else {
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
		}
		return
	}

	app.render(w, http.StatusOK, "home.html", data)
}

func (app *Application) aboutView(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "about.html", data)
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
