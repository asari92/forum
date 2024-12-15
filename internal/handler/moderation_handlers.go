package handler

import (
	"net/http"

	"forum/pkg/validator"
)

func (app *Application) moderationUnapprovedPostsView(w http.ResponseWriter, r *http.Request) {
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

	paginationURL := "/moderation/posts/unapproved"
	unapprovedPostsDTO, err := app.Service.Post.GetAllPaginatedUnapprovedPostsDTO(page, pageSize, paginationURL)
	if err != nil {
		app.Logger.Error("get unapproved posts", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	data := app.newTemplateData(r)
	data.Posts = unapprovedPostsDTO.Posts
	data.Header = "Unapproved posts"
	data.Pagination = pagination{
		CurrentPage:      unapprovedPostsDTO.CurrentPage,
		HasNextPage:      unapprovedPostsDTO.HasNextPage,
		PaginationAction: unapprovedPostsDTO.PaginationURL,
	}
	app.render(w, http.StatusOK, "user_posts.html", data)
}
