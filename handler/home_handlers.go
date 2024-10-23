package handler

import (
	"fmt"
	"forum/entities"
	"net/http"
	"strconv"
	"strings"
)

type pagination struct {
	CurrentPage      int
	HasNextPage      bool
	PaginationAction string
}

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	// Определяем текущую страницу. По умолчанию - страница 1.
	page := 1
	pageSize := 10 // Количество постов на одной странице

	// Получаем посты для нужной страницы.
	posts, err := app.Service.Post.GetAllPaginatedPosts(page, pageSize)
	if err != nil {
		app.Logger.Error("get all paginated posts", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	// Проверяем, есть ли следующая страница
	hasNextPage := len(posts) > pageSize

	// Если постов больше, чем pageSize, обрезаем список до pageSize
	if hasNextPage {
		posts = posts[:pageSize]
	}

	categories, err := app.Service.Category.GetAll()
	if err != nil {
		app.Logger.Error("get all categories", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	form := postCreateForm{
		Categories: []int{},
	}

	data := app.newTemplateData(r)
	data.Posts = posts
	data.Categories = categories
	data.Header = "All posts"
	data.Pagination = pagination{
		CurrentPage:      page,
		HasNextPage:      hasNextPage,
		PaginationAction: "/", // Маршрут для пагинации
	}
	data.Form = form
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

	form := postCreateForm{
		Categories: categoryIDs,
	}

	allCategories, err := app.Service.Category.GetAll()
	if err != nil {
		app.Logger.Error("get all categories", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	// Проверка категорий
	form.validateCategories(allCategories)

	data := app.newTemplateData(r)
	data.Categories = allCategories
	data.Form = form
	data.Header = "All posts"

	if !form.Valid() {
		posts, err := app.Service.Post.GetAllPaginatedPosts(page, pageSize)
		if err != nil {
			app.Logger.Error("get all posts", "error", err)
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
			return
		}

		// Проверяем, есть ли следующая страница
		hasNextPage := len(posts) > pageSize

		// Если постов больше, чем pageSize, обрезаем список до pageSize
		if hasNextPage {
			posts = posts[:pageSize]
		}

		data.Posts = posts
		data.Pagination = pagination{
			CurrentPage:      page,
			HasNextPage:      hasNextPage,
			PaginationAction: "/", // Маршрут для пагинации
		}

		app.render(w, http.StatusUnprocessableEntity, "home.html", data)
		return
	}

	// Получаем посты с пагинацией и на одну запись больше
	posts, err := app.Service.Post.GetPaginatedPostsByCategory(form.Categories, page, pageSize)
	if err != nil {
		app.Logger.Error("get paginated posts by category", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	// Проверяем, есть ли следующая страница
	hasNextPage := len(posts) > pageSize

	// Если постов больше, чем pageSize, обрезаем список до pageSize
	if hasNextPage {
		posts = posts[:pageSize]
	}

	// Собираем названия категорий, совпадающие с выбранными categoryIDs
	var filteredCategoryNames []string
	for _, category := range allCategories {
		for _, selectedID := range categoryIDs {
			if category.ID == selectedID {
				filteredCategoryNames = append(filteredCategoryNames, category.Name)
			}
		}
	}

	// Если категории выбраны, создаем строку заголовка с их названиями
	if len(filteredCategoryNames) > 0 {
		data.Header = fmt.Sprintf("Posts filtered by: %s", strings.Join(filteredCategoryNames, ", "))
	}

	data.Posts = posts
	data.Pagination = pagination{
		CurrentPage:      page,
		HasNextPage:      hasNextPage,
		PaginationAction: "/", // Маршрут для пагинации
	}

	app.render(w, http.StatusOK, "home.html", data)
}

func (form *postCreateForm) validateCategories(allCategories []*entities.Category) {
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
		form.AddFieldError("categories", "Need one or more category")
	}
}

func (app *Application) aboutView(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "about.html", data)
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
