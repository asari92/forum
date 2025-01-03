package handler

import (
	"errors"
	"fmt"
	"net/http"

	"forum/internal/entities"
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

func (app *Application) moderationApprovePost(w http.ResponseWriter, r *http.Request) {
	sess := app.SessionFromContext(r)
	userID, ok := sess.Get(AuthUserIDSessionKey).(int)
	if !ok || userID < 1 {
		err := errors.New("get userID in accountPasswordUpdate")
		app.Logger.Error("get userid from session", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return
	}

	postID, err := validator.ValidateID(r.PathValue("post_id"))
	if err != nil {
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return
	}

	err = app.Service.Post.ApprovePost(postID)
	if err != nil {
		app.Logger.Error("post approval", "error", err)
		if errors.Is(err, entities.ErrNoRecord) {
			app.render(w, http.StatusBadRequest, Errorpage, nil)
		} else {
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
		}
		return
	}

	err = sess.Set(FlashSessionKey, "Post successfully aprroved!")
	if err != nil {
		// кажется тут не нужна ошибка, достаточно логирования
		app.Logger.Error("Session error during set flash", "error", err)
	}

	http.Redirect(w, r, "/moderation/posts/unapproved", http.StatusSeeOther)
}

func (app *Application) moderationReportPost(w http.ResponseWriter, r *http.Request) {
	sess := app.SessionFromContext(r)
	userID, ok := sess.Get(AuthUserIDSessionKey).(int)
	if !ok || userID < 1 {
		err := errors.New("get userID in accountPasswordUpdate")
		app.Logger.Error("get userid from session", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return
	}

	postID, err := validator.ValidateID(r.PathValue("post_id"))
	if err != nil {
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return
	}

	reason := r.PostForm.Get("report_reason")
	if !validator.Matches(reason, validator.TextRX) {
		app.render(w, http.StatusBadRequest, Errorpage,
			&templateData{AppError: AppError{Message: "Report reason must contain only english or russian letters", StatusCode: http.StatusBadRequest}})
		return
	}

	err = app.Service.Reaction.CreateReport(userID, postID, reason)
	if err != nil {
		app.Logger.Error("create report", "error", err)
		if errors.Is(err, entities.ErrNoRecord) {
			app.render(w, http.StatusBadRequest, Errorpage, nil)
		} else {
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
		}
		return
	}

	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

func (app *Application) moderationApplicationView(w http.ResponseWriter, r *http.Request) {
	sess := app.SessionFromContext(r)
	userId, ok := sess.Get(AuthUserIDSessionKey).(int)
	if !ok || userId < 1 {
		err := errors.New("get userID in editPostView")
		app.Logger.Error("get userid from session", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	userRole, ok := sess.Get(UserRoleSessionKey).(string)
	if !ok {
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return

	}
	data := app.newTemplateData(r)

	if userRole == entities.RoleUser {
		app.render(w, http.StatusOK, "moderation_application.html", data)
	}
}

func (app *Application) createModerationApplication(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return
	}
	sess := app.SessionFromContext(r)
	userId, ok := sess.Get(AuthUserIDSessionKey).(int)
	if !ok || userId < 1 {
		err := errors.New("get userID in editPostView")
		app.Logger.Error("get userid from session", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	userRole, ok := sess.Get(UserRoleSessionKey).(string)
	if !ok {
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return

	}

	form := app.Service.User.NewModerationForm()
	form.Reason = r.PostForm.Get("reason")

	if userRole == entities.RoleUser {
		err := app.Service.User.CreateModerationRequest(userId, &form)
		if err != nil {
			if errors.Is(err, entities.ErrInvalidData) || errors.Is(err, entities.ErrFormAlreadySubmitted) {
				data := app.newTemplateData(r)
				data.Form = form
				app.render(w, http.StatusUnprocessableEntity, "moderation_application.html", data)
			} else {
				app.render(w, http.StatusInternalServerError, Errorpage, nil)
			}
			return

		}

	}

	http.Redirect(w, r, "/account/view", http.StatusSeeOther)
}

func (app *Application) moderatorsView(w http.ResponseWriter, r *http.Request) {
	sess := app.SessionFromContext(r)
	userId, ok := sess.Get(AuthUserIDSessionKey).(int)
	if !ok || userId < 1 {
		err := errors.New("get userID in editPostView")
		app.Logger.Error("get userid from session", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	userRole, ok := sess.Get(UserRoleSessionKey).(string)
	if !ok {
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return

	}
	moderators, err := app.Service.User.GetModerators()
	if err != nil {
		if errors.Is(err, entities.ErrNoRecord) {
		} else {

			app.render(w, http.StatusInternalServerError, Errorpage, nil)
			return

		}
	}
	data := app.newTemplateData(r)
	data.Users = moderators

	if userRole == entities.RoleAdmin {
		app.render(w, http.StatusOK, "moderatorslist.html", data)
	}
}

func (app *Application) deleteModerator(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return
	}
	sess := app.SessionFromContext(r)
	userId, ok := sess.Get(AuthUserIDSessionKey).(int)
	if !ok || userId < 1 {
		err := errors.New("get userID in editPostView")
		app.Logger.Error("get userid from session", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	userRole, ok := sess.Get(UserRoleSessionKey).(string)
	if !ok {
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return

	}

	moderatorId, err := validator.ValidateID(r.PostForm.Get("id"))
	if err != nil {
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return
	}

	if userRole == entities.RoleAdmin {
		err := app.Service.User.DeleteModerator(moderatorId)
		if err != nil {
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
			return

		}

	}

	http.Redirect(w, r, "/moderators/list", http.StatusSeeOther)
}

func (app *Application) moderationApplicantsView(w http.ResponseWriter, r *http.Request) {
	sess := app.SessionFromContext(r)
	userId, ok := sess.Get(AuthUserIDSessionKey).(int)
	if !ok || userId < 1 {
		err := errors.New("get userID in editPostView")
		app.Logger.Error("get userid from session", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	userRole, ok := sess.Get(UserRoleSessionKey).(string)
	if !ok {
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return

	}
	applicants, err := app.Service.User.GetModerationApplicants()
	if err != nil {
		if errors.Is(err, entities.ErrNoRecord) {
		} else {

			app.render(w, http.StatusInternalServerError, Errorpage, nil)
			return

		}
	}
	data := app.newTemplateData(r)
	data.Applicants = applicants

	if userRole == entities.RoleAdmin {
		app.render(w, http.StatusOK, "moderation-applicants.html", data)
	}
}

func (app *Application) requestModeratorRole(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return
	}
	sess := app.SessionFromContext(r)
	userId, ok := sess.Get(AuthUserIDSessionKey).(int)
	if !ok || userId < 1 {
		err := errors.New("get userID in editPostView")
		app.Logger.Error("get userid from session", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	userRole, ok := sess.Get(UserRoleSessionKey).(string)
	if !ok {
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return

	}

	applicantId, err := validator.ValidateID(r.PostForm.Get("id"))
	if err != nil {
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return
	}

	if userRole == entities.RoleAdmin {
		err := app.Service.User.ApproveModerationRequest(applicantId)
		if err != nil {
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
			return
		}
		err = app.Service.User.DeleteModerationRequest(applicantId)
		if err != nil {
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
			return
		}

	}

	http.Redirect(w, r, "/moderation-applicants", http.StatusSeeOther)
}

func (app *Application) rejectModeratorRequest(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return
	}
	sess := app.SessionFromContext(r)
	userId, ok := sess.Get(AuthUserIDSessionKey).(int)
	if !ok || userId < 1 {
		err := errors.New("get userID in editPostView")
		app.Logger.Error("get userid from session", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	userRole, ok := sess.Get(UserRoleSessionKey).(string)
	if !ok {
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return

	}

	applicantId, err := validator.ValidateID(r.PostForm.Get("id"))
	if err != nil {
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return
	}

	if userRole == entities.RoleAdmin {
		err := app.Service.User.DeleteModerationRequest(applicantId)
		if err != nil {
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
			return
		}

	}

	http.Redirect(w, r, "/moderation-applicants", http.StatusSeeOther)
}

func (app *Application) acceptReport(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return
	}
	sess := app.SessionFromContext(r)
	userId, ok := sess.Get(AuthUserIDSessionKey).(int)
	if !ok || userId < 1 {
		err := errors.New("get userID in editPostView")
		app.Logger.Error("get userid from session", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	userRole, ok := sess.Get(UserRoleSessionKey).(string)
	if !ok {
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return

	}

	postId, err := validator.ValidateID(r.PostForm.Get("postId"))
	if err != nil {
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return
	}

	if userRole == entities.RoleAdmin {
		err := app.Service.Post.DeletePost(postId, userId)
		if err != nil {
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
			return
		}

	}

	http.Redirect(w, r, "/administration/reports", http.StatusSeeOther)
}

func (app *Application) rejectReport(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return
	}
	sess := app.SessionFromContext(r)
	userId, ok := sess.Get(AuthUserIDSessionKey).(int)
	if !ok || userId < 1 {
		err := errors.New("get userID in editPostView")
		app.Logger.Error("get userid from session", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	userRole, ok := sess.Get(UserRoleSessionKey).(string)
	if !ok {
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return

	}

	postId, err := validator.ValidateID(r.PostForm.Get("postId"))
	if err != nil {
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return
	}
	reporterId, err := validator.ValidateID(r.PostForm.Get("reporter_id"))
	if err != nil {
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return
	}

	if userRole == entities.RoleAdmin {
		err := app.Service.Post.DeleteReport(reporterId, postId)
		if err != nil {
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
			return
		}

	}

	http.Redirect(w, r, "/administration/reports", http.StatusSeeOther)
}

func (app *Application) categoryEditView(w http.ResponseWriter, r *http.Request) {
	sess := app.SessionFromContext(r)
	userId, ok := sess.Get(AuthUserIDSessionKey).(int)
	if !ok || userId < 1 {
		err := errors.New("get userID in editPostView")
		app.Logger.Error("get userid from session", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	userRole, ok := sess.Get(UserRoleSessionKey).(string)
	if !ok {
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return

	}

	categories, err := app.Service.Category.GetAll() //дальше больше надо будет передавать
	if err != nil {
		app.Logger.Error("get all categories", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	data := app.newTemplateData(r)
	data.Categories = categories

	if userRole == entities.RoleAdmin {
		app.render(w, http.StatusOK, "categoryedit.html", data)
	}
}

func (app *Application) createCategory(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return
	}
	sess := app.SessionFromContext(r)
	userId, ok := sess.Get(AuthUserIDSessionKey).(int)
	if !ok || userId < 1 {
		err := errors.New("get userID in editPostView")
		app.Logger.Error("get userid from session", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	userRole, ok := sess.Get(UserRoleSessionKey).(string)
	if !ok {
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return

	}

	form := app.Service.Category.NewCategoryCreateForm()
	form.Name = r.PostForm.Get("category_name")

	if userRole == entities.RoleAdmin {
		_, err := app.Service.Category.Insert(&form)
		if err != nil {
			if errors.Is(err, entities.ErrInvalidData) {
				data := app.newTemplateData(r)
				categories, err := app.Service.Category.GetAll() //дальше больше надо будет передавать
				if err != nil {
					app.Logger.Error("get all categories", "error", err)
					app.render(w, http.StatusInternalServerError, Errorpage, nil)
					return
				}
				data.Categories = categories

				data.Form = form
				app.render(w, http.StatusUnprocessableEntity, "categoryedit.html", data)
			} else {
				app.render(w, http.StatusInternalServerError, Errorpage, nil)
			}
			return

		}

	}

	http.Redirect(w, r, "/edit/category", http.StatusSeeOther)
}

func (app *Application) deleteCategory(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return
	}
	sess := app.SessionFromContext(r)
	userId, ok := sess.Get(AuthUserIDSessionKey).(int)
	if !ok || userId < 1 {
		err := errors.New("get userID in editPostView")
		app.Logger.Error("get userid from session", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	userRole, ok := sess.Get(UserRoleSessionKey).(string)
	if !ok {
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return

	}

	categoryId, err := validator.ValidateID(r.PostForm.Get("category_id"))
	if err != nil {
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return

	}
	if userRole == entities.RoleAdmin {
		err := app.Service.Category.Delete(categoryId)
		fmt.Println(err)
		if err != nil {
			if errors.Is(err, entities.ErrInvalidData) || errors.Is(err, entities.ErrNoRecord) {
				data := app.newTemplateData(r)
				categories, err := app.Service.Category.GetAll() //дальше больше надо будет передавать
				if err != nil {
					app.Logger.Error("get all categories", "error", err)
					app.render(w, http.StatusInternalServerError, Errorpage, nil)
					return
				}
				data.Categories = categories
				app.render(w, http.StatusUnprocessableEntity, "categoryedit.html", data)
			} else {
				app.render(w, http.StatusInternalServerError, Errorpage, nil)
			}
			return

		}

	}

	http.Redirect(w, r, "/edit/category", http.StatusSeeOther)
}
