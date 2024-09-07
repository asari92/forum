package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"time"

	"forum/internal/session"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	// Initialize a new buffer.
	buf := new(bytes.Buffer)

	// Write the template to the buffer, instead of straight to the
	// http.ResponseWriter. If there's an error, call our serverError() helper
	// and then return.
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// If the template is written to the buffer without any errors, we are safe
	// to go ahead and write the HTTP status code to http.ResponseWriter.
	w.WriteHeader(status)

	// Write the contents of the buffer to the http.ResponseWriter. Note: this
	// is another time where we pass our http.ResponseWriter to a function that
	// takes an io.Writer.
	buf.WriteTo(w)
}

func (app *application) SessionFromContext(r *http.Request) session.Session {
	return r.Context().Value("session").(session.Session)
}

func (app *application) newTemplateData(w http.ResponseWriter, r *http.Request) *templateData {
	sess := app.SessionFromContext(r)
	flash, ok := sess.Get("flash").(string)
	if ok {
		err := sess.Delete("flash")
		if err != nil {
			app.serverError(w, err)
		}
	}
	return &templateData{
		CurrentYear:     time.Now().Year(),
		Flash:           flash,
		CSRFToken:       r.Context().Value("csrfToken").(string),
		IsAuthenticated: app.isAuthenticated(r),
	}
}

func (app *application) isAuthenticated(r *http.Request) bool {
	sess := app.SessionFromContext(r)
	userId := sess.Get("authenticatedUserID")
	return userId != nil
}

// Генерация CSRF-токена
func (app *application) generateCSRFToken() string {
	h := md5.New()
	salt := "YouShallNotPass%^7&8888" // Используйте уникальную соль для своего приложения
	io.WriteString(h, salt+time.Now().String())
	token := fmt.Sprintf("%x", h.Sum(nil))
	return token
}
