package handler

import (
	"errors"
	"net/http"
	"strings"

	"forum/internal/entities"
)

func (app *Application) userSignupView(w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(r.Referer(), "/user/signup") && !strings.Contains(r.Referer(), "/user/login") {
		sess := app.SessionFromContext(r)
		sess.Set(RedirectPathAfterLoginSessionKey, r.Referer())
	}

	data := app.newTemplateData(r)
	data.Form = app.Service.User.NewUserAuthForm()
	app.render(w, http.StatusOK, "signup.html", data)
}

func (app *Application) userSignup(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return
	}

	form := app.Service.User.NewUserAuthForm()
	form.Username = r.PostForm.Get("username")
	form.Email = r.PostForm.Get("email")
	form.Password = r.PostForm.Get("password")

	err = app.Service.User.Insert(&form)
	if err != nil {
		if errors.Is(err, entities.ErrInvalidCredentials) {
			form.AddFieldError("insert user credentials", "invalid credentials")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
		} else if errors.Is(err, entities.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
		} else {
			app.Logger.Error("insert user credentials", "error", err)
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
		}
		return
	}

	sess := app.SessionFromContext(r)
	// Если валидация прошла успешно, удаляем токен из сессии
	err = sess.Delete(CsrfTokenSessionKey)
	if err != nil {
		app.Logger.Error("Session error during delete csrfToken", "error", err)
	}

	err = sess.Set(FlashSessionKey, "Your signup was successful. Please log in.")
	if err != nil {
		app.Logger.Error("set flashsessionkey", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *Application) userLoginView(w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(r.Referer(), "/user/signup") && !strings.Contains(r.Referer(), "/user/login") {
		sess := app.SessionFromContext(r)
		sess.Set(RedirectPathAfterLoginSessionKey, r.Referer())
	}

	data := app.newTemplateData(r)
	data.Form = app.Service.User.NewUserAuthForm()

	app.render(w, http.StatusOK, "login.html", data)
}

func (app *Application) userLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return
	}

	form := app.Service.User.NewUserAuthForm()
	form.Email = r.PostForm.Get("email")
	form.Password = r.PostForm.Get("password")

	// Check whether the credentials are valid. If they're not, add a generic
	// non-field error message and re-display the login page.
	userID, err := app.Service.User.Authenticate(&form)
	if err != nil {
		if errors.Is(err, entities.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.html", data)
		} else {
			app.Logger.Error("get id Authenticate user", "error", err)
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
			return
		}
		return
	}

	sess := app.SessionFromContext(r)
	// Если валидация прошла успешно, удаляем токен из сессии
	err = sess.Delete(CsrfTokenSessionKey)
	if err != nil {
		app.Logger.Error("Session error during delete csrfToken", "error", err)
	}

	// Add the ID of the current user to the session, so that they are now
	// 'logged in'.
	err = sess.Set(AuthUserIDSessionKey, userID)
	if err != nil {
		app.Logger.Error("set AuthUserIDSessionKey", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	redirectUrl := "/"
	path, ok := sess.Get(RedirectPathAfterLoginSessionKey).(string)
	if ok {
		err = sess.Delete(RedirectPathAfterLoginSessionKey)
		if err != nil {
			app.Logger.Error("Session error during delete redirectPath", "error", err)
		}
		redirectUrl = path
	}

	err = sess.Set(FlashSessionKey, "Your log in was successful.")
	if err != nil {
		app.Logger.Error("Set FlashSessionKey", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	err = app.SessionManager.RenewToken(w, r, userID)
	if err != nil {
		app.Logger.Error("renewtoken", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
}

func (app *Application) userLogout(w http.ResponseWriter, r *http.Request) {
	sess := app.SessionFromContext(r)
	err := sess.Delete(AuthUserIDSessionKey)
	if err != nil {
		app.Logger.Error("Session error during delete authUserID", "error", err)
	}
	sess.Set(FlashSessionKey, "You've been logged out successfully!")

	// err = app.SessionManager.RenewToken(w, r)
	// if err != nil {
	// 	app.Logger.Error("renewtoken", "error", err)
	// 	app.render(w, http.StatusInternalServerError, Errorpage, nil)
	// 	return
	// }

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
