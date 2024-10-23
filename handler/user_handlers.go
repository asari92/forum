package handler

import (
	"errors"
	"forum/internal/validator"
	"forum/repository"
	"net/http"
	"strings"
)

type signupForm struct {
	Username string
	Email    string
	Password string
	validator.Validator
}

type userLoginForm struct {
	Email    string
	Password string
	validator.Validator
}

func (app *Application) userSignupView(w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(r.Referer(), "/user/signup") && !strings.Contains(r.Referer(), "/user/login") {
		sess := app.SessionFromContext(r)
		sess.Set(RedirectPathAfterLoginSessionKey, r.Referer())
	}

	data := app.newTemplateData(r)
	data.Form = signupForm{}
	app.render(w, http.StatusOK, "signup.html", data)
}

func (app *Application) userSignup(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return
	}
	form := signupForm{
		Username: r.PostForm.Get("username"),
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}

	// ДОБАВИТЬ ПРОВЕРОК!!!!!!!!!!!!!!!!!!!
	form.CheckField(validator.NotBlank(form.Username), "username", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Username, 100), "username", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.MaxChars(form.Email, 100), "email", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")
	form.CheckField(validator.MaxChars(form.Password, 100), "password", "This field cannot be more than 100 characters long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
		return
	}

	err = app.Service.User.Insert(form.Username, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateEmail) {
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
	data.Form = userLoginForm{}

	app.render(w, http.StatusOK, "login.html", data)
}

func (app *Application) userLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return
	}
	form := userLoginForm{
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.MaxChars(form.Email, 100), "email", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")
	form.CheckField(validator.MaxChars(form.Password, 100), "password", "This field cannot be more than 100 characters long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.html", data)
		return
	}

	// Check whether the credentials are valid. If they're not, add a generic
	// non-field error message and re-display the login page.
	id, err := app.Service.User.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, repository.ErrInvalidCredentials) {
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

	err = sess.Set(FlashSessionKey, "Your log in was successful.")
	if err != nil {
		app.Logger.Error("Set FlashSessionKey", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	// Add the ID of the current user to the session, so that they are now
	// 'logged in'.
	err = sess.Set(AuthUserIDSessionKey, id)
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

	err = app.SessionManager.RenewToken(w, r)
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

	err = app.SessionManager.RenewToken(w, r)
	if err != nil {
		app.Logger.Error("renewtoken", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
