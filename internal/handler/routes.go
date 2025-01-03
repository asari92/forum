package handler

import (
	"net/http"
	"os"

	"forum/ui"
)

// Кастомная файловая система, которая запрещает доступ к директориям
type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if s.IsDir() {
		return nil, os.ErrNotExist
	}

	return f, nil
}

func (app *Application) Routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(neuteredFileSystem{http.FS(ui.Files)})
	mux.Handle("GET /static/", cacheControlMiddleware(fileServer))

	uploadServer := http.FileServer(http.Dir("./uploads"))
	mux.Handle("GET /uploads/", http.StripPrefix("/uploads/", cacheControlMiddleware(uploadServer)))

	dynamic := New(app.verifyCSRF, app.sessionMiddleware, app.authenticate)
	mux.Handle("/", dynamic.ThenFunc(app.errorHandler))
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("POST /", dynamic.ThenFunc(app.filterPosts))
	mux.Handle("GET /post/view/{id}", dynamic.ThenFunc(app.postView))
	mux.Handle("GET /commented-post/view/{id}", dynamic.ThenFunc(app.commentedPostView))
	mux.Handle("GET /user/{userId}/posts", dynamic.ThenFunc(app.userPostsView))
	mux.Handle("POST /user/{userId}/posts", dynamic.ThenFunc(app.userPostsView))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignupView))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLoginView))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("GET /about", dynamic.ThenFunc(app.aboutView))

	mux.Handle("GET /auth/google/login", dynamic.ThenFunc(app.oauthGoogleLogin))
	mux.Handle("GET /auth/google/callback", dynamic.ThenFunc(app.oauthGoogleCallback))
	mux.Handle("GET /auth/github/login", dynamic.ThenFunc(app.oauthGithubLogin))
	mux.Handle("GET /auth/github/callback", dynamic.ThenFunc(app.oauthGithubCallback))

	protected := dynamic.Append(app.requireAuthentication)
	mux.Handle("GET /post/edit/{post_id}", protected.ThenFunc(app.editPostView))
	mux.Handle("POST /post/edit/{post_id}", protected.ThenFunc(app.editPost))
	mux.Handle("POST /post/view/{id}", protected.ThenFunc(app.postReaction))
	mux.Handle("POST /post/delete", protected.ThenFunc(app.DeletePost))
	mux.Handle("GET /post/create", protected.ThenFunc(app.postCreateView))
	mux.Handle("POST /post/create", protected.ThenFunc(app.postCreate))

	mux.Handle("GET /comment/edit", protected.ThenFunc(app.editCommentView))
	mux.Handle("POST /comment/edit", protected.ThenFunc(app.editComment))
	mux.Handle("POST /comment/delete", protected.ThenFunc(app.DeleteComment))

	mux.Handle("GET /account/notification", protected.ThenFunc(app.notificationView))
	mux.Handle("GET /account/view", protected.ThenFunc(app.accountView))
	mux.Handle("GET /account/password/update", protected.ThenFunc(app.accountPasswordUpdateView))
	mux.Handle("POST /account/password/update", protected.ThenFunc(app.accountPasswordUpdate))

	mux.Handle("GET /user/liked", protected.ThenFunc(app.userLikedPostsView))
	mux.Handle("GET /user/commented", protected.ThenFunc(app.userCommentedPostsView))
	mux.Handle("POST /user/liked", protected.ThenFunc(app.userLikedPostsView))
	mux.Handle("POST /user/commented", protected.ThenFunc(app.userCommentedPostsView))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogout))

	mux.Handle("GET /moderation-application", protected.ThenFunc(app.moderationApplicationView))
	mux.Handle("POST /moderation-application", protected.ThenFunc(app.createModerationApplication))

	moderated := protected.Append(app.requireModeration)
	mux.Handle("GET /moderation/posts/unapproved", moderated.ThenFunc(app.moderationUnapprovedPostsView))
	mux.Handle("POST /moderation/approve/{post_id}", protected.ThenFunc(app.moderationApprovePost))
	mux.Handle("POST /moderation/report/{post_id}", protected.ThenFunc(app.moderationReportPost))


	administrated := protected.Append(app.requireAdministration)
	mux.Handle("GET /administration/reports", administrated.ThenFunc(app.administrationReportsView))
	mux.Handle("GET /moderation-applicants", administrated.ThenFunc(app.moderationApplicantsView))
	mux.Handle("GET /moderators/list", administrated.ThenFunc(app.moderatorsView))
	mux.Handle("POST /moderators/delete", administrated.ThenFunc(app.deleteModerator))
	mux.Handle("POST /moderation/accept", administrated.ThenFunc(app.requestModeratorRole))
	mux.Handle("POST /moderation/reject", administrated.ThenFunc(app.rejectModeratorRequest))
	mux.Handle("POST /report/accept", administrated.ThenFunc(app.acceptReport))
	mux.Handle("POST /report/reject", administrated.ThenFunc(app.rejectReport))
	mux.Handle("GET /edit/category", administrated.ThenFunc(app.categoryEditView))
	mux.Handle("POST /admin/category/create", administrated.ThenFunc(app.createCategory))
	mux.Handle("POST /admin/category/delete", administrated.ThenFunc(app.deleteCategory))


	


	standard := New(app.recoverPanic, app.logRequest, secureHeaders, app.rateLimiting)
	return standard.Then(mux)
}
