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
	mux.Handle("GET /static/", fileServer)

	uploadServer := http.FileServer(http.Dir("./uploads"))
	mux.Handle("GET /uploads/", http.StripPrefix("/uploads/", uploadServer))

	dynamic := New(app.verifyCSRF, app.sessionMiddleware, app.authenticate)

	mux.Handle("/", dynamic.ThenFunc(app.errorHandler))
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("POST /", dynamic.ThenFunc(app.filterPosts))
	mux.Handle("GET /post/view/{id}", dynamic.ThenFunc(app.postView))
	mux.Handle("GET /user/{userId}/posts", dynamic.ThenFunc(app.userPostsView))
	mux.Handle("POST /user/{userId}/posts", dynamic.ThenFunc(app.userPostsView))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignupView))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLoginView))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("GET /about", dynamic.ThenFunc(app.aboutView))

	protected := dynamic.Append(app.requireAuthentication)

	mux.Handle("POST /post/view/{id}", protected.ThenFunc(app.postReaction))
	mux.Handle("GET /post/create", protected.ThenFunc(app.postCreateView))
	mux.Handle("POST /post/create", protected.ThenFunc(app.postCreate))
	mux.Handle("GET /account/view", protected.ThenFunc(app.accountView))
	mux.Handle("GET /user/liked", protected.ThenFunc(app.userLikedPostsView))
	mux.Handle("POST /user/liked", protected.ThenFunc(app.userLikedPostsView))
	mux.Handle("GET /account/password/update", protected.ThenFunc(app.accountPasswordUpdateView))
	mux.Handle("POST /account/password/update", protected.ThenFunc(app.accountPasswordUpdate))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogout))

	standard := New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(mux)
}
