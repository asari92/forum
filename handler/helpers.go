package handler

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"time"

	"forum/entities"
	"forum/internal/session"
)

var RequestPaths = map[string]bool{
	"/":                        true,
	"/ping":                    true,
	"/about":                   true,
	"/user/login":              true,
	"/user/signup":             true,
	"/post/create":             true,
	"/account/view":            true,
	"/user/liked":              true,
	"/account/password/update": true,
	"/user/logout":             true,
}

const Errorpage = "errorpage.html"
const DefaultCategory = 1

// func (app *application) serverErrorLogging(err error) {
// 	_, path, line, _ := runtime.Caller(1)

// 	app.logger.Error("Server error occurred",
// 		"error", err,
// 		path, line,
// 		// "stack_trace", string(debug.Stack()), // Включение стека
// 	)
// }

func (app *Application) serverError(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *Application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := app.TemplateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.Logger.Error("Server error occured", "render error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	// Initialize a new buffer.
	buf := new(bytes.Buffer)

	if page == Errorpage {
		data := &templateData{
			AppError: AppError{StatusCode: status, Message: http.StatusText(status)},
		}
		err := ts.Execute(buf, data)
		if err != nil {
			app.Logger.Error("Server error occured", "render error", err)
			app.serverError(w)
			return

		}

	} else {
		// Write the template to the buffer, instead of straight to the
		// http.ResponseWriter. If there's an error, call our serverError() helper
		// and then return.
		err := ts.ExecuteTemplate(buf, "base", data)
		if err != nil {
			app.Logger.Error("Server error occured", "render error", err)
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
			return
		}
	}

	// If the template is written to the buffer without any errors, we are safe
	// to go ahead and write the HTTP status code to http.ResponseWriter.
	w.WriteHeader(status)

	// Write the contents of the buffer to the http.ResponseWriter. Note: this
	// is another time where we pass our http.ResponseWriter to a function that
	// takes an io.Writer.
	buf.WriteTo(w)
}

func (app *Application) SessionFromContext(r *http.Request) session.Session {
	return r.Context().Value(sessionContextKey).(session.Session)
}

func (app *Application) newTemplateData(r *http.Request) *templateData {
	sess := app.SessionFromContext(r)
	flash, ok := sess.Get(FlashSessionKey).(string)
	if ok {
		err := sess.Delete(FlashSessionKey)
		if err != nil {
			app.Logger.Error("Session error during delete flash", "error", err)
		}
	}
	return &templateData{
		CurrentYear:     time.Now().Year(),
		Flash:           flash,
		CSRFToken:       r.Context().Value(csrfTokenContextKey).(string),
		IsAuthenticated: app.isAuthenticated(r),
		ReactionData:    &ReactionData{UserReaction: &entities.PostReaction{}},
	}
}

func (app *Application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}

	return isAuthenticated
}

// Генерация CSRF-токена
func (app *Application) generateCSRFToken() string {
	h := md5.New()
	salt := "YouShallNotPass%^7&8888" // Используйте уникальную соль для своего приложения
	io.WriteString(h, salt+time.Now().String())
	token := fmt.Sprintf("%x", h.Sum(nil))
	return token
}
