package handler

import (
	"errors"
	"fmt"
	"net/http"

	"forum/internal/entities"
	"forum/pkg/validator"
)

func (app *Application) postView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return
	}

	postID, err := validator.ValidateID(r.PathValue("id"))
	if err != nil {
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return
	}

	sess := app.SessionFromContext(r)
	userID, ok := sess.Get(AuthUserIDSessionKey).(int)
	if !ok {
		userID = 0
	}

	postData, err := app.Service.Post.GetPostDTO(postID, userID)
	if err != nil {
		app.Logger.Error("get post data", "error", err)
		if errors.Is(err, entities.ErrNoRecord) {
			app.render(w, http.StatusBadRequest, Errorpage, nil)
		} else {
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
		}
		return
	}

	data := app.newTemplateData(r)
	data.User = &entities.User{}
	data.User.ID = userID
	data.Post = postData.Post
	data.Comments = postData.Comments
	data.Categories = postData.Categories
	data.Images = postData.Images
	data.ReactionData.Likes = postData.Likes
	data.ReactionData.Dislikes = postData.Dislikes
	data.ReactionData.UserReaction = postData.UserReaction
	data.Form = sess.Get(ReactionFormSessionKey)
	err = sess.Delete(ReactionFormSessionKey)
	if err != nil {
		app.Logger.Error("Session error during delete reaction form", "error", err)
	}

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
	form := app.Service.Post.NewPostCreateForm()
	// первая категория всегда отмечена
	form.Categories = []int{DefaultCategory}
	data.Form = form

	app.render(w, http.StatusOK, "create_post.html", data)
}

func (app *Application) DeletePost(w http.ResponseWriter, r *http.Request) {
	sess := app.SessionFromContext(r)
	userId, ok := sess.Get(AuthUserIDSessionKey).(int)
	if !ok || userId < 1 {
		err := errors.New("get userID in postCreate")
		app.Logger.Error("get userid from session", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}
	err := r.ParseForm()
	if err != nil {
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return
	}
	post_id, err := validator.ValidateID(r.PostForm.Get("post_id"))
	if err != nil {
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return
	}

	err = app.Service.DeletePost(post_id, userId)
	if err != nil {

		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Application) postCreate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(20*1024*1024 + (10 * 1024))
	if err != nil {
		app.Logger.Error("Uploaded file is too big", "error", err)
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return

	}

	sess := app.SessionFromContext(r)
	userId, ok := sess.Get(AuthUserIDSessionKey).(int)
	if !ok || userId < 1 {
		err := errors.New("get userID in postCreate")
		app.Logger.Error("get userid from session", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	var categoryIDs []int
	for _, id := range r.PostForm["categories"] {
		intID, err := validator.ValidateID(id)
		if err != nil {
			app.render(w, http.StatusBadRequest, Errorpage, nil)
			return
		}
		categoryIDs = append(categoryIDs, intID)
	}

	form := app.Service.Post.NewPostCreateForm()
	form.Title = r.PostForm.Get("title")
	form.Content = r.PostForm.Get("content")
	form.Categories = categoryIDs
	files := r.MultipartForm.File["image"]

	postId, allCategories, err := app.Service.Post.CreatePostWithCategories(&form, files, userId)
	if err != nil {
		app.Logger.Error("insert post and categories", "error", err)
		if errors.Is(err, entities.ErrInvalidCredentials) {
			data := app.newTemplateData(r)
			data.Categories = allCategories
			// первая категория будет отмечена по умолчанию
			if len(form.Categories) == 0 {
				form.Categories = []int{DefaultCategory}
			}
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "create_post.html", data)
		} else if errors.Is(err, entities.ErrNoRecord) {
			app.render(w, http.StatusBadRequest, Errorpage, nil)
		} else {
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
		}
		return
	}

	err = sess.Set(FlashSessionKey, "Post successfully created!")
	if err != nil {
		// кажется тут не нужна ошибка, достаточно логирования
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

	userId, err := validator.ValidateID(r.PathValue("userId"))
	if err != nil || userId < 1 {
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return
	}

	page := 1
	pageSize := 10

	if p, err := validator.ValidateID(r.PostFormValue("page")); err == nil {
		page = p
	}

	paginationURL := fmt.Sprintf("/user/%d/posts", userId)
	app.Logger.Debug("get user posts", "userID", userId, "page", page, "pageSize", pageSize, "paginationURL", paginationURL)
	userPostsDTO, err := app.Service.Post.GetUserPostsDTO(userId, page, pageSize, paginationURL)
	app.Logger.Debug("get user posts", "userPostsDTO", userPostsDTO)
	if err != nil {
		app.Logger.Error("get user posts", "error", err)
		if errors.Is(err, entities.ErrNoRecord) {
			app.render(w, http.StatusBadRequest, Errorpage, nil)
		} else {
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Posts = userPostsDTO.Posts
	data.Header = fmt.Sprintf("Posts by %s", userPostsDTO.User.Username)
	data.Pagination = pagination{
		CurrentPage:      userPostsDTO.CurrentPage,
		HasNextPage:      userPostsDTO.HasNextPage,
		PaginationAction: userPostsDTO.PaginationURL,
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

	if p, err := validator.ValidateID(r.PostFormValue("page")); err == nil {
		page = p
	}

	paginationURL := "/user/liked"
	userLikedPostsDTO, err := app.Service.Post.GetUserLikedPostsDTO(userID, page, pageSize, paginationURL)
	if err != nil {
		app.Logger.Error("get user liked posts", "error", err)
		if errors.Is(err, entities.ErrNoRecord) {
			app.render(w, http.StatusBadRequest, Errorpage, nil)
		} else {
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Posts = userLikedPostsDTO.Posts
	data.Header = fmt.Sprintf("Posts liked by %s", userLikedPostsDTO.User.Username)
	data.Pagination = pagination{
		CurrentPage:      userLikedPostsDTO.CurrentPage,
		HasNextPage:      userLikedPostsDTO.HasNextPage,
		PaginationAction: userLikedPostsDTO.PaginationURL,
	}
	app.render(w, http.StatusOK, "user_posts.html", data)
}

func (app *Application) userCommentedPostsView(w http.ResponseWriter, r *http.Request) {
	sess := app.SessionFromContext(r)
	userID, ok := sess.Get(AuthUserIDSessionKey).(int)
	if !ok || userID < 1 {
		err := errors.New("get userID in userCommentedPostsView")
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

	if p, err := validator.ValidateID(r.PostFormValue("page")); err == nil {
		page = p
	}

	paginationURL := "/user/commented"
	userCommentedPostsDTO, err := app.Service.Post.GetUserCommentedPostsDTO(userID, page, pageSize, paginationURL)
	if err != nil {
		app.Logger.Error("get user commented posts", "error", err)
		if errors.Is(err, entities.ErrNoRecord) {
			app.render(w, http.StatusBadRequest, Errorpage, nil)
		} else {
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Posts = userCommentedPostsDTO.Posts
	data.Header = fmt.Sprintf("Posts commented by %s", userCommentedPostsDTO.User.Username)
	data.Pagination = pagination{
		CurrentPage:      userCommentedPostsDTO.CurrentPage,
		HasNextPage:      userCommentedPostsDTO.HasNextPage,
		PaginationAction: userCommentedPostsDTO.PaginationURL,
	}
	app.render(w, http.StatusOK, "user_posts.html", data)
}
