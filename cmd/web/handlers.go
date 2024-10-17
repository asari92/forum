package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"forum/internal/models"
	"forum/internal/validator"
)

func (app *application) errorHandler(w http.ResponseWriter, r *http.Request) {
	if RequestPaths[r.URL.Path] {
		app.render(w, http.StatusMethodNotAllowed, Errorpage, nil)
		return
	} 

	app.render(w, http.StatusNotFound, Errorpage, nil)
}

type pagination struct {
	CurrentPage      int
	HasNextPage      bool
	PaginationAction string
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Определяем текущую страницу. По умолчанию - страница 1.
	page := 1
	pageSize := 10 // Количество постов на одной странице

	// Получаем посты для нужной страницы.
	posts, err := app.posts.GetAllPaginatedPosts(page, pageSize)
	if err != nil {
		app.logger.Error("get all paginated posts", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	// Проверяем, есть ли следующая страница
	hasNextPage := len(posts) > pageSize

	// Если постов больше, чем pageSize, обрезаем список до pageSize
	if hasNextPage {
		posts = posts[:pageSize]
	}

	categories, err := app.categories.GetAll()
	if err != nil {
		app.logger.Error("get all categories", "error", err)
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

func (app *application) filterPosts(w http.ResponseWriter, r *http.Request) {
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

	allCategories, err := app.categories.GetAll()
	if err != nil {
		app.logger.Error("get all categories", "error", err)
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
		posts, err := app.posts.GetAllPaginatedPosts(page, pageSize)
		if err != nil {
			app.logger.Error("get all posts", "error", err)
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
	posts, err := app.posts.GetPaginatedPostsByCategory(form.Categories, page, pageSize)
	if err != nil {
		app.logger.Error("get paginated posts by category", "error", err)
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

func (app *application) postView(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || postID < 1 {
		app.render(w, http.StatusNotFound, Errorpage, nil)
		return
	}

	post, err := app.posts.Get(postID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.render(w, http.StatusNotFound, Errorpage, nil)
			return
		} else {
			app.logger.Error("get post from database", "error", err)
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
			return
		}
	}

	categories, err := app.categories.GetCategoriesForPost(postID)
	if err != nil {
		app.logger.Error("get categories for post", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	likes, dislikes, err := app.postReactions.GetReactionsCount(postID)
	if err != nil {
		app.logger.Error("get likes and dislikes count for post", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	sess := app.SessionFromContext(r)
	var userReaction *models.PostReaction
	userID, ok := sess.Get(AuthUserIDSessionKey).(int)
	if ok && userID != 0 {
		userReaction, err = app.postReactions.GetUserReaction(userID, postID) // Получите реакцию пользователя
		if err != nil {
			app.logger.Error("get post user reaction", "error", err)
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
			return
		}
	}
	comments, err := app.comments.GetComments(postID)
	if err != nil {
		app.logger.Error("get comments from database", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}
	// continue
	for _, val := range comments {
		userID, ok := sess.Get(AuthUserIDSessionKey).(int)
		if ok && userID != 0 {
			userReaction, err := app.commentReactions.GetUserReaction(userID, val.ID)
			if err != nil {
				app.logger.Error("get comment user reaction", "error", err)
				app.render(w, http.StatusInternalServerError, Errorpage, nil)
				return

			}
			if userReaction != nil {
				if userReaction.IsLike {
					val.UserReaction = 1
				} else {
					val.UserReaction = -1
				}
			}

		}

		like, err := app.commentReactions.GetLikesCount(val.ID)
		if err != nil {
			app.logger.Error("get comment likes count", "error", err)
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
			return
		}
		val.Like = like
		app.logger.Debug("comment likes count", "like:", like)
		dislike, err := app.commentReactions.GetDislikesCount(val.ID)
		if err != nil {
			app.logger.Error("get comment dislike count", "error", err)
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
			return
		}
		val.Dislike = dislike

	}

	data := app.newTemplateData(r)
	data.Post = post
	data.Categories = categories
	data.Comments = comments
	data.ReactionData.Likes = likes
	data.ReactionData.Dislikes = dislikes
	data.ReactionData.UserReaction = userReaction

	app.render(w, http.StatusOK, "post_view.html", data)
}

func (app *application) postReaction(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return
	}

	postID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || postID < 1 {
		app.render(w, http.StatusNotFound, Errorpage, nil)
		return
	}

	comment := r.PostForm.Get("comment_content")
	postIsLike := r.PostForm.Get("post_is_like")
	commentIsLike := r.PostForm.Get("comment_is_like")

	sess := app.SessionFromContext(r)
	userID, ok := sess.Get(AuthUserIDSessionKey).(int)
	if !ok || userID < 1 {
		err = errors.New("get userID in postReaction")
		app.logger.Error("get userid from session", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	if comment != "" {
		err = app.comments.InsertComment(postID, userID, comment)
		if err != nil {
			app.logger.Error("insert comment to database", "error", err)
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
			return

		}
	} else if postIsLike != "" {
		// Преобразуем isLike в bool
		like := postIsLike == "true"

		var userReaction *models.PostReaction
		userReaction, err = app.postReactions.GetUserReaction(userID, postID) // Получите реакцию пользователя
		if err != nil {
			app.logger.Error("get post user reaction", "error", err)
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
			return
		}

		if userReaction != nil && userReaction.IsLike == like {
			err = app.postReactions.RemoveReaction(userID, postID)
			if err != nil {
				app.logger.Error("remove post reaction", "error", err)
				app.render(w, http.StatusInternalServerError, Errorpage, nil)
				return
			}
		} else {
			err = app.postReactions.AddReaction(userID, postID, like)
			if err != nil {
				app.logger.Error("add post reaction", "error", err)
				app.render(w, http.StatusInternalServerError, Errorpage, nil)
				return
			}
		}
	} else if commentIsLike != "" {
		commentID, err := strconv.Atoi(r.PostForm.Get("comment_id"))
		if err != nil || commentID < 1 {
			app.render(w, http.StatusBadRequest, Errorpage, nil)
			return
		}

		reaction := commentIsLike == "true"
		var commentReaction *models.CommentReaction

		commentReaction, err = app.commentReactions.GetUserReaction(userID, commentID)
		if err != nil {
			app.logger.Error("get user comment reaction", "error", err)
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
			return
		}

		if commentReaction != nil && reaction == commentReaction.IsLike {
			err = app.commentReactions.RemoveReaction(userID, commentID)
			if err != nil {
				app.logger.Error("remove comment reaction", "error", err)
				app.render(w, http.StatusInternalServerError, Errorpage, nil)
				return
			}
		} else {
			err = app.commentReactions.AddReaction(userID, commentID, reaction)
			if err != nil {

				app.logger.Error("add comment reaction", "error", err)
				app.render(w, http.StatusInternalServerError, Errorpage, nil)
				return
			}
		}

	}

	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

// /comment/reaction
func (app *application) commentReaction(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return
	}

	commentID, err := strconv.Atoi(r.PostForm.Get("comment_id"))
	if err != nil || commentID < 1 {
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return
	}
	isLike := r.PostForm.Get("is_like")
	reaction := isLike == "true"
	var commentReaction *models.CommentReaction
	sess := app.SessionFromContext(r)
	userID, ok := sess.Get(AuthUserIDSessionKey).(int)
	if ok && userID != 0 {
		commentReaction, err = app.commentReactions.GetUserReaction(userID, commentID)
		if err != nil {
			app.logger.Error("get user reaction", "error", err)
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
			return
		}
	}

	if reaction == commentReaction.IsLike {
		err = app.commentReactions.RemoveReaction(userID, commentID)
		if err != nil {
			app.logger.Error("remove reaction", "error", err)
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
			return
		} else {
			err = app.commentReactions.AddReaction(userID, commentID, reaction)
			if err != nil {

				app.logger.Error("Add Reaction", "error", err)
				app.render(w, http.StatusInternalServerError, Errorpage, nil)
				return
			}
		}
	}

	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

func (app *application) userPostsView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return
	}

	userId, err := strconv.Atoi(r.PathValue("userId"))
	if err != nil || userId < 1 {
		app.render(w, http.StatusNotFound, Errorpage, nil)
		return
	}

	page := 1
	pageSize := 10

	if p, err := strconv.Atoi(r.PostFormValue("page")); err == nil && p > 0 {
		page = p
	}

	posts, err := app.posts.GetUserPaginatedPosts(userId, page, pageSize)
	if err != nil {
		app.logger.Error("get user posts", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	user, err := app.users.Get(userId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			if len(posts) == 0 {
				app.render(w, http.StatusNotFound, Errorpage, nil)
				return
			} else {
				user = &models.User{Username: "Deleted User", ID: userId}
			}
		} else {
			app.logger.Error("get user query error", "error", err)
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
			return
		}
	}

	hasNextPage := len(posts) > pageSize

	if hasNextPage {
		posts = posts[:pageSize]
	}

	data := app.newTemplateData(r)
	data.Posts = posts
	data.Header = fmt.Sprintf("Posts by %s", user.Username)
	data.Pagination = pagination{
		CurrentPage:      page,
		HasNextPage:      hasNextPage,
		PaginationAction: fmt.Sprintf("/user/%d/posts", userId), // Маршрут для пагинации
	}
	app.render(w, http.StatusOK, "user_posts.html", data)
}

func (app *application) userLikedPostsView(w http.ResponseWriter, r *http.Request) {
	sess := app.SessionFromContext(r)
	userID, ok := sess.Get(AuthUserIDSessionKey).(int)
	if !ok || userID < 1 {
		err := errors.New("get userID in userLikedPostsView")
		app.logger.Error("get userid", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
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

	posts, err := app.posts.GetUserLikedPaginatedPosts(userID, page, pageSize)
	if err != nil {
		app.logger.Error("get user liked posts", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	user, err := app.users.Get(userID)
	if err != nil {
		app.logger.Error("get user", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	hasNextPage := len(posts) > pageSize

	if hasNextPage {
		posts = posts[:pageSize]
	}

	data := app.newTemplateData(r)
	data.Posts = posts
	data.Header = fmt.Sprintf("Posts liked by %s", user.Username)
	data.Pagination = pagination{
		CurrentPage:      page,
		HasNextPage:      hasNextPage,
		PaginationAction: "/user/liked", // Маршрут для пагинации
	}
	app.render(w, http.StatusOK, "user_posts.html", data)
}

func (app *application) postCreateView(w http.ResponseWriter, r *http.Request) {
	categories, err := app.categories.GetAll()
	if err != nil {
		app.logger.Error("get all categories", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
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
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return
	}

	var categoryIDs []int
	for _, id := range r.PostForm["categories"] {
		intID, err := strconv.Atoi(id)
		if err != nil {
			app.render(w, http.StatusBadRequest, Errorpage, nil)
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
		app.logger.Error("get all categories", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	form.validateCategories(allCategories)

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Categories = allCategories
		// первая категория будет отмечена по умолчанию
		if len(form.Categories) == 0 {
			form.Categories = []int{models.DefaultCategory}
		}
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create_post.html", data)
		return
	}

	sess := app.SessionFromContext(r)
	userId, ok := sess.Get(AuthUserIDSessionKey).(int)
	if !ok || userId < 1 {
		err = errors.New("get userID in postCreate")
		app.logger.Error("get userid from session", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	postId, err := app.posts.InsertPostWithCategories(form.Title, form.Content, userId, form.Categories)
	if err != nil {
		app.logger.Error("insert post and categories", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	err = sess.Set(FlashSessionKey, "Post successfully created!")
	if err != nil {
		// кажется тут не нужна ошибка, достаточно логирования
		// app.serverError(w, err)
		// trace := fmt.Sprintf("Error: %+v", err) // Log full stack trace
		// app.errorLog.Println(trace)
		app.logger.Error("Session error during set flash", "error", err)
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
		form.AddFieldError("categories", "Need one or more category")
	}
}

func (app *application) userSignupView(w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(r.Referer(), "/user/signup") && !strings.Contains(r.Referer(), "/user/login") {
		sess := app.SessionFromContext(r)
		sess.Set(RedirectPathAfterLoginSessionKey, r.Referer())
	}

	data := app.newTemplateData(r)
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
		app.render(w, http.StatusBadRequest, Errorpage, nil)
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
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
		return
	}

	err = app.users.Insert(form.Username, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
		} else {
			app.logger.Error("insert user credentials", "error", err)
			app.render(w, http.StatusInternalServerError, Errorpage, nil)

		}
		return
	}
	sess := app.SessionFromContext(r)
	// Если валидация прошла успешно, удаляем токен из сессии
	err = sess.Delete(CsrfTokenSessionKey)
	if err != nil {
		app.logger.Error("Session error during delete csrfToken", "error", err)
	}

	err = sess.Set(FlashSessionKey, "Your signup was successful. Please log in.")
	if err != nil {
		app.logger.Error("set flashsessionkey", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

type userLoginForm struct {
	Email    string
	Password string
	validator.Validator
}

func (app *application) userLoginView(w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(r.Referer(), "/user/signup") && !strings.Contains(r.Referer(), "/user/login") {
		sess := app.SessionFromContext(r)
		sess.Set(RedirectPathAfterLoginSessionKey, r.Referer())
	}

	data := app.newTemplateData(r)
	data.Form = userLoginForm{}

	app.render(w, http.StatusOK, "login.html", data)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return
	}
	form := userLoginForm{
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.MaxChars(form.Email, 100), "email", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")
	form.CheckField(validator.MaxChars(form.Password, 100), "password", "This field cannot be more than 100 characters long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.html", data)
		return
	}

	// Check whether the credentials are valid. If they're not, add a generic
	// non-field error message and re-display the login page.
	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.html", data)
		} else {
			app.logger.Error("get id Authenticate user", "error", err)
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
			return
		}
		return
	}

	sess := app.SessionFromContext(r)
	// Если валидация прошла успешно, удаляем токен из сессии
	err = sess.Delete(CsrfTokenSessionKey)
	if err != nil {
		app.logger.Error("Session error during delete csrfToken", "error", err)
	}

	err = sess.Set(FlashSessionKey, "Your log in was successful.")
	if err != nil {
		app.logger.Error("Set FlashSessionKey", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	// Add the ID of the current user to the session, so that they are now
	// 'logged in'.
	err = sess.Set(AuthUserIDSessionKey, id)
	if err != nil {
		app.logger.Error("set AuthUserIDSessionKey", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	redirectUrl := "/"
	path, ok := sess.Get(RedirectPathAfterLoginSessionKey).(string)
	if ok {
		err = sess.Delete(RedirectPathAfterLoginSessionKey)
		if err != nil {
			app.logger.Error("Session error during delete redirectPath", "error", err)
		}
		redirectUrl = path
	}

	err = app.sessionManager.RenewToken(w, r)
	if err != nil {
		app.logger.Error("renewtoken", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
}

func (app *application) userLogout(w http.ResponseWriter, r *http.Request) {
	sess := app.SessionFromContext(r)
	err := sess.Delete(AuthUserIDSessionKey)
	if err != nil {
		app.logger.Error("Session error during delete authUserID", "error", err)
	}
	sess.Set(FlashSessionKey, "You've been logged out successfully!")

	err = app.sessionManager.RenewToken(w, r)
	if err != nil {
		app.logger.Error("renewtoken", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) aboutView(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "about.html", data)
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (app *application) accountView(w http.ResponseWriter, r *http.Request) {
	sess := app.SessionFromContext(r)
	userID, ok := sess.Get(AuthUserIDSessionKey).(int)
	if !ok || userID < 1 {
		err := errors.New("get userID in accountView")
		app.logger.Error("get userid from session", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	user, err := app.users.Get(userID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		} else {
			app.logger.Error("get user from db", "error", err)
			app.render(w, http.StatusInternalServerError, Errorpage, nil)

		}
		return
	}

	data := app.newTemplateData(r)
	data.User = user

	app.render(w, http.StatusOK, "account.html", data)
}

type accountPasswordUpdateForm struct {
	CurrentPassword         string
	NewPassword             string
	NewPasswordConfirmation string
	validator.Validator
}

func (app *application) accountPasswordUpdateView(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = accountPasswordUpdateForm{}

	app.render(w, http.StatusOK, "password.html", data)
}

func (app *application) accountPasswordUpdate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return
	}
	form := accountPasswordUpdateForm{
		CurrentPassword:         r.PostForm.Get("currentPassword"),
		NewPassword:             r.PostForm.Get("newPassword"),
		NewPasswordConfirmation: r.PostForm.Get("newPasswordConfirmation"),
	}

	form.CheckField(validator.NotBlank(form.CurrentPassword), "currentPassword", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.NewPassword), "newPassword", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.NewPassword, 8), "newPassword", "This field must be at least 8 characters long")
	form.CheckField(validator.NotBlank(form.NewPasswordConfirmation), "newPasswordConfirmation", "This field cannot be blank")
	form.CheckField(form.NewPassword == form.NewPasswordConfirmation, "newPasswordConfirmation", "Passwords do not match")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form

		app.render(w, http.StatusUnprocessableEntity, "password.html", data)
		return
	}

	sess := app.SessionFromContext(r)
	userID, ok := sess.Get(AuthUserIDSessionKey).(int)
	if !ok || userID < 1 {
		err = errors.New("get userID in accountPasswordUpdate")
		app.logger.Error("get userid from session", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	err = app.users.PasswordUpdate(userID, form.CurrentPassword, form.NewPassword)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddFieldError("currentPassword", "Current password is incorrect")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "password.html", data)
		} else {
			app.logger.Error("password update", "error", err)
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
		}
		return
	}

	sess.Set("flash", "Your password has been updated!")

	http.Redirect(w, r, "/account/view", http.StatusSeeOther)
}
