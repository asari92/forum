package handler

import (
	"errors"
	"net/http"

	"forum/internal/entities"
)

func (app *Application) accountView(w http.ResponseWriter, r *http.Request) {
	sess := app.SessionFromContext(r)
	userID, ok := sess.Get(AuthUserIDSessionKey).(int)
	if !ok || userID < 1 {
		err := errors.New("get userID in accountView")
		app.Logger.Error("get userid from session", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	user, err := app.Service.User.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, entities.ErrNoRecord) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		} else {
			app.Logger.Error("get user from db", "error", err)
			app.render(w, http.StatusInternalServerError, Errorpage, nil)

		}
		return
	}

	data := app.newTemplateData(r)
	data.User = user

	app.render(w, http.StatusOK, "account.html", data)
}

func (app *Application) accountPasswordUpdateView(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = app.Service.User.NewAccountPasswordUpdateForm()

	app.render(w, http.StatusOK, "password.html", data)
}

func (app *Application) accountPasswordUpdate(w http.ResponseWriter, r *http.Request) {
	sess := app.SessionFromContext(r)
	userID, ok := sess.Get(AuthUserIDSessionKey).(int)
	if !ok || userID < 1 {
		err := errors.New("get userID in accountPasswordUpdate")
		app.Logger.Error("get userid from session", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return
	}

	form := app.Service.User.NewAccountPasswordUpdateForm()
	form.CurrentPassword = r.PostForm.Get("currentPassword")
	form.NewPassword = r.PostForm.Get("newPassword")
	form.NewPasswordConfirmation = r.PostForm.Get("newPasswordConfirmation")

	err = app.Service.User.UpdatePassword(userID, &form)
	if err != nil {
		if errors.Is(err, entities.ErrInvalidCredentials) {
			form.AddFieldError("invalid credentials", "invalid credentials")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "password.html", data)
		} else {
			app.Logger.Error("password update", "error", err)
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
		}
		return
	}

	sess.Set("flash", "Your password has been updated!")

	http.Redirect(w, r, "/account/view", http.StatusSeeOther)
}