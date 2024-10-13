package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"forum/internal/models"
	"forum/internal/session"
	"io"
	"net/http"
	"runtime"
	"time"
)

func (app *application) renderErrorPage(w http.ResponseWriter, statuscode int) {

	data := &templateData{
		AppError: AppError{StatusCode: statuscode, Message: http.StatusText(statuscode)},
	}

	ts, ok := app.templateCache["errorpage.html"]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", "errorpage.html")
		app.serverErrorLogging(err)
		app.serverError(w)
		return
	}
	buf := new(bytes.Buffer)

	err := ts.Execute(buf, data)
	if err != nil {
		app.serverErrorLogging(err)
		app.serverError(w)

		return
	}

	w.WriteHeader(statuscode)

	if _, err := buf.WriteTo(w); err != nil {
		app.serverErrorLogging(err)
		app.serverError(w)

		return
	}
}

func (app *application) serverErrorLogging(err error) {
	_, path, line, _ := runtime.Caller(1)

	app.logger.Error("Server error occurred",
		"error", err,
		path, line,
		// "stack_trace", string(debug.Stack()), // Включение стека
	)
}

func (app *application) serverError(w http.ResponseWriter) {

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}



func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) error {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)

		return err
	}

	// Initialize a new buffer.
	buf := new(bytes.Buffer)

	// Write the template to the buffer, instead of straight to the
	// http.ResponseWriter. If there's an error, call our serverError() helper
	// and then return.
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {

		return err
	}

	// If the template is written to the buffer without any errors, we are safe
	// to go ahead and write the HTTP status code to http.ResponseWriter.
	w.WriteHeader(status)

	// Write the contents of the buffer to the http.ResponseWriter. Note: this
	// is another time where we pass our http.ResponseWriter to a function that
	// takes an io.Writer.
	buf.WriteTo(w)
	return nil
}

func (app *application) SessionFromContext(r *http.Request) session.Session {
	return r.Context().Value(sessionContextKey).(session.Session)
}

func (app *application) newTemplateData(r *http.Request) *templateData {
	sess := app.SessionFromContext(r)
	flash, ok := sess.Get(FlashSessionKey).(string)
	if ok {
		err := sess.Delete(FlashSessionKey)
		if err != nil {
			app.logger.Error("Session error during delete flash", "error", err)
		}
	}
	return &templateData{
		CurrentYear:     time.Now().Year(),
		Flash:           flash,
		CSRFToken:       r.Context().Value(csrfTokenContextKey).(string),
		IsAuthenticated: app.isAuthenticated(r),
		ReactionData:    &ReactionData{UserReaction: &models.PostReaction{}},
	}
}

func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}

	return isAuthenticated
}

// Генерация CSRF-токена
func (app *application) generateCSRFToken() string {
	h := md5.New()
	salt := "YouShallNotPass%^7&8888" // Используйте уникальную соль для своего приложения
	io.WriteString(h, salt+time.Now().String())
	token := fmt.Sprintf("%x", h.Sum(nil))
	return token
}
