package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	dynamic := New(app.verifyCSRF, app.sessionMiddleware, app.authenticate)

	mux.Handle("GET /", dynamic.ThenFunc(app.errorHandler))
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /post/view/{id}", dynamic.ThenFunc(app.postView))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignupView))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLoginView))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLogin))

	protected := dynamic.Append(app.requireAuthentication)

	mux.Handle("GET /post/create", protected.ThenFunc(app.postCreateView))
	mux.Handle("POST /post/create", protected.ThenFunc(app.postCreate))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogout))

	standard := New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(mux)
}
