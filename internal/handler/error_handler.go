package handler

import "net/http"

func (app *Application) errorHandler(w http.ResponseWriter, r *http.Request) {
	if RequestPaths[r.URL.Path] {
		app.render(w, http.StatusMethodNotAllowed, Errorpage, nil)
		return
	}

	app.render(w, http.StatusNotFound, Errorpage, nil)
}
