package handler

import (
	"errors"
	"forum/internal/validator"
	"forum/repository"
	"net/http"
)

type accountPasswordUpdateForm struct {
	CurrentPassword         string
	NewPassword             string
	NewPasswordConfirmation string
	validator.Validator
}

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
		if errors.Is(err, repository.ErrNoRecord) {
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
	data.Form = accountPasswordUpdateForm{}

	app.render(w, http.StatusOK, "password.html", data)
}

func (app *Application) accountPasswordUpdate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.render(w, http.StatusBadRequest, Errorpage, nil)
		return
	}
	form := accountPasswordUpdateForm{
		CurrentPassword:         r.PostForm.Get("currentPassword"),
		NewPassword:             r.PostForm.Get("newPassword"),
		NewPasswordConfirmation: r.PostForm.Get("newPasswordConfirmation"),
	}

	form.CheckField(validator.NotBlank(form.CurrentPassword), "currentPassword", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.NewPassword), "newPassword", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.NewPassword, 8), "newPassword", "This field must be at least 8 characters long")
	form.CheckField(validator.NotBlank(form.NewPasswordConfirmation), "newPasswordConfirmation", "This field cannot be blank")
	form.CheckField(form.NewPassword == form.NewPasswordConfirmation, "newPasswordConfirmation", "Passwords do not match")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form

		app.render(w, http.StatusUnprocessableEntity, "password.html", data)
		return
	}

	sess := app.SessionFromContext(r)
	userID, ok := sess.Get(AuthUserIDSessionKey).(int)
	if !ok || userID < 1 {
		err = errors.New("get userID in accountPasswordUpdate")
		app.Logger.Error("get userid from session", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	err = app.Service.User.UpdatePassword(userID, form.CurrentPassword, form.NewPassword)
	if err != nil {
		if errors.Is(err, repository.ErrInvalidCredentials) {
			form.AddFieldError("currentPassword", "Current password is incorrect")
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