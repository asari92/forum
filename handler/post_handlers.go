package handler

import (
	"errors"
	"fmt"
	"forum/entities"
	"forum/internal/validator"
	"forum/repository"
	"net/http"
	"strconv"
)

type postCreateForm struct {
	Title      string
	Content    string
	Categories []int
	validator.Validator
}

func (app *Application) postView(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || postID < 1 {
		app.render(w, http.StatusNotFound, Errorpage, nil)
		return
	}

	post, err := app.Service.Post.GetPost(postID)
	if err != nil {
		if errors.Is(err, repository.ErrNoRecord) {
			app.render(w, http.StatusNotFound, Errorpage, nil)
			return
		} else {
			app.Logger.Error("get post from database", "error", err)
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
			return
		}
	}

	categories, err := app.Service.Category.GetCategoriesForPost(postID)
	if err != nil {
		app.Logger.Error("get categories for post", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	likes, dislikes, err := app.Service.PostReaction.GetReactionsCount(postID)
	if err != nil {
		app.Logger.Error("get likes and dislikes count for post", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	sess := app.SessionFromContext(r)
	var userReaction *entities.PostReaction
	userID, ok := sess.Get(AuthUserIDSessionKey).(int)
	if ok && userID != 0 {
		userReaction, err = app.Service.PostReaction.GetUserReaction(userID, postID) // Получите реакцию пользователя
		if err != nil {
			app.Logger.Error("get post user reaction", "error", err)
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
			return
		}
	}
	comments, err := app.Service.Comment.GetPostComments(postID)
	if err != nil {
		app.Logger.Error("get comments from database", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}
	// continue
	for _, val := range comments {
		userID, ok := sess.Get(AuthUserIDSessionKey).(int)
		if ok && userID != 0 {
			userReaction, err := app.Service.CommentReaction.GetUserReaction(userID, val.ID)
			if err != nil {
				app.Logger.Error("get comment user reaction", "error", err)
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

		like, err := app.Service.CommentReaction.GetLikesCount(val.ID)
		if err != nil {
			app.Logger.Error("get comment likes count", "error", err)
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
			return
		}
		val.Like = like
		app.Logger.Debug("comment likes count", "like:", like)
		dislike, err := app.Service.CommentReaction.GetDislikesCount(val.ID)
		if err != nil {
			app.Logger.Error("get comment dislike count", "error", err)
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

func (app *Application) postCreateView(w http.ResponseWriter, r *http.Request) {
	categories, err := app.Service.Category.GetAll()
	if err != nil {
		app.Logger.Error("get all categories", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	data := app.newTemplateData(r)
	data.Categories = categories
	data.Form = postCreateForm{
		// первая категория всегда отмечена
		Categories: []int{DefaultCategory},
	}

	app.render(w, http.StatusOK, "create_post.html", data)
}

func (app *Application) postCreate(w http.ResponseWriter, r *http.Request) {
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

	allCategories, err := app.Service.Category.GetAll()
	if err != nil {
		app.Logger.Error("get all categories", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	form.validateCategories(allCategories)

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Categories = allCategories
		// первая категория будет отмечена по умолчанию
		if len(form.Categories) == 0 {
			form.Categories = []int{DefaultCategory}
		}
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create_post.html", data)
		return
	}

	sess := app.SessionFromContext(r)
	userId, ok := sess.Get(AuthUserIDSessionKey).(int)
	if !ok || userId < 1 {
		err = errors.New("get userID in postCreate")
		app.Logger.Error("get userid from session", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	postId, err := app.Service.Post.CreatePostWithCategories(form.Title, form.Content, userId, form.Categories)
	if err != nil {
		app.Logger.Error("insert post and categories", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	err = sess.Set(FlashSessionKey, "Post successfully created!")
	if err != nil {
		// кажется тут не нужна ошибка, достаточно логирования
		// app.serverError(w, err)
		// trace := fmt.Sprintf("Error: %+v", err) // Log full stack trace
		// app.errorLog.Println(trace)
		app.Logger.Error("Session error during set flash", "error", err)
	}

	http.Redirect(w, r, fmt.Sprintf("/post/view/%d", postId), http.StatusSeeOther)
}

func (app *Application) userPostsView(w http.ResponseWriter, r *http.Request) {
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

	posts, err := app.Service.Post.GetUserPaginatedPosts(userId, page, pageSize)
	if err != nil {
		app.Logger.Error("get user posts", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	user, err := app.Service.User.GetUserByID(userId)
	if err != nil {
		if errors.Is(err, repository.ErrNoRecord) {
			if len(posts) == 0 {
				app.render(w, http.StatusNotFound, Errorpage, nil)
				return
			} else {
				user = &entities.User{Username: "Deleted User", ID: userId}
			}
		} else {
			app.Logger.Error("get user query error", "error", err)
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

func (app *Application) userLikedPostsView(w http.ResponseWriter, r *http.Request) {
	sess := app.SessionFromContext(r)
	userID, ok := sess.Get(AuthUserIDSessionKey).(int)
	if !ok || userID < 1 {
		err := errors.New("get userID in userLikedPostsView")
		app.Logger.Error("get userid", "error", err)
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

	posts, err := app.Service.Post.GetUserLikedPaginatedPosts(userID, page, pageSize)
	if err != nil {
		app.Logger.Error("get user liked posts", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	user, err := app.Service.User.GetUserByID(userID)
	if err != nil {
		app.Logger.Error("get user", "error", err)
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
