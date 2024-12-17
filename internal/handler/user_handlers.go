package handler

import (
	"errors"
	"net/http"
	"strings"

	"forum/internal/entities"
	"forum/pkg/validator"
)

func (app *Application) notificationView(w http.ResponseWriter, r *http.Request) {
	sess := app.SessionFromContext(r)
	userId, ok := sess.Get(AuthUserIDSessionKey).(int)
	if !ok || userId < 1 {
		err := errors.New("get userID in editPostView")
		app.Logger.Error("get userid from session", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}
	notifications, err := app.Service.Post.GetUserNotifications(userId)
	if err != nil {
		if errors.Is(err, entities.ErrNoRecord) {
			app.render(w, http.StatusBadRequest, Errorpage, nil)
		} else {
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Notifications = notifications

	data.Form = sess.Get(ReactionFormSessionKey)
	err = sess.Delete(ReactionFormSessionKey)
	if err != nil {
		app.Logger.Error("Session error during delete reaction form", "error", err)
	}

	app.render(w, http.StatusOK, "notification.html", data)
}

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

	form.CheckField(validator.Matches(form.Username, validator.UsernameRX), "username", "This field must be a valid username")
	form.CheckField(validator.NotBlank(form.Username), "username", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Username, 100), "username", "This field cannot be more than 100 characters long")

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.MaxChars(form.Email, 100), "email", "This field cannot be more than 100 characters long")

	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.NotHaveAnySpaces(form.Password), "password", "This field cannot have any spaces")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")
	form.CheckField(validator.MaxChars(form.Password, 100), "password", "This field cannot be more than 100 characters long")
	form.CheckField(validator.Matches(form.Password, validator.PasswordRX), "password", "This field must contain only letters, digits and these symbols !_.@#$%^&*")

	if !form.Valid() {
		form.AddFieldError("insert user credentials", "invalid credentials")
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
		return
	}

	_, err = app.Service.User.Insert(form.Username, form.Email, form.Password, entities.RoleUser)
	if err != nil {
		if errors.Is(err, entities.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
		} else if errors.Is(err, entities.ErrDuplicateUsername) {
			form.AddFieldError("username", "Username is already in use")
		} else {
			app.Logger.Error("insert user credentials", "error", err)
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
			return
		}
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
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
		// app.render(w, http.StatusInternalServerError, Errorpage, nil)
		// return
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

	// Валидация данных
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.MaxChars(form.Email, 100), "email", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")
	form.CheckField(validator.MaxChars(form.Password, 100), "password", "This field cannot be more than 100 characters long")

	if !form.Valid() {
		form.AddNonFieldError("Email or password is incorrect")
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.html", data)
		return
	}

	// Check whether the credentials are valid. If they're not, add a generic
	// non-field error message and re-display the login page.
	user, err := app.Service.User.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, entities.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.html", data)
			return
		} else {
			app.Logger.Error("get id Authenticate user", "error", err)
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
			return
		}
	}

	sess := app.SessionFromContext(r)
	// Если валидация прошла успешно, удаляем токен из сессии
	err = sess.Delete(CsrfTokenSessionKey)
	if err != nil {
		app.Logger.Error("Session error during delete csrfToken", "error", err)
	}

	// Add the ID of the current user to the session, so that they are now
	// 'logged in'.
	err = sess.Set(AuthUserIDSessionKey, user.ID)
	if err != nil {
		app.Logger.Error("set AuthUserIDSessionKey", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	err = sess.Set(UserRoleSessionKey, user.Role)
	if err != nil {
		app.Logger.Error("set UserRoleSessionKey", "error", err)
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
		// app.render(w, http.StatusInternalServerError, Errorpage, nil)
		// return
	}

	err = app.SessionManager.RenewToken(w, r, user.ID)
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

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
