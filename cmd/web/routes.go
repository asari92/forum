package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /", app.errorHandler)
	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /post/view/{id}", app.postView)
	mux.HandleFunc("GET /post/create", app.postCreateView)
	mux.HandleFunc("POST /post/create", app.postCreate)

	myChain := New(app.recoverPanic, app.logRequest)
	myOtherChain := myChain.Append(secureHeaders)

	return myOtherChain.Then(mux)
}
