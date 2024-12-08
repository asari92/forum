package handler

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"

	"forum/internal/entities"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Scopes: OAuth 2.0 scopes provide a way to limit the amount of access that is granted to an access token.
var googleOauthConfig = &oauth2.Config{
	RedirectURL:  "",
	ClientID:     "",
	ClientSecret: "",
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	},
	Endpoint: google.Endpoint,
}

var oauthState string = "state-token"

func (app *Application) oauthGoogleLogin(w http.ResponseWriter, r *http.Request) {
	oauthState = generateStateOauthCookie(w)
	googleOauthConfig.ClientID = app.Config.GoogleClientID
	googleOauthConfig.ClientSecret = app.Config.GoogleClientSecret
	googleOauthConfig.RedirectURL = app.Config.GoogleClientCallbackURL
	url := googleOauthConfig.AuthCodeURL(oauthState, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	return state
}

// Callback от Google
func (app *Application) oauthGoogleCallback(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != oauthState {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	code := r.FormValue("code")
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	client := googleOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userInfo struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Failed to parse user info", http.StatusInternalServerError)
		return
	}

	// Логика регистрации/входа пользователя
	// log.Printf("User info: %s %s (%s)", userInfo.Name, userInfo.Email)
	// w.Write([]byte("Welcome, " + userInfo.Name))

	userID, err := app.Service.User.OauthAuthenticate(userInfo.Email)
	if err != nil {
		if errors.Is(err, entities.ErrInvalidCredentials) {
		} else {
			app.Logger.Error("get id Authenticate user", "error", err)
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
			return
		}
	}

	sess := app.SessionFromContext(r)

	if userID < 1 {
		userID, err = app.Service.User.Insert(userInfo.Name, userInfo.Email, "")
		if err != nil {
			if errors.Is(err, entities.ErrDuplicateUsername) {
				err = sess.Set(FlashSessionKey, "Sorry, but this username is already in use, you must register with email and password or change name in your google account.")
				if err != nil {
					app.Logger.Error("set flashsessionkey", "error", err)
					app.render(w, http.StatusInternalServerError, Errorpage, nil)
					return
				}
				data := app.newTemplateData(r)
				app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
				return
			} else {
				app.Logger.Error("insert user credentials", "error", err)
				app.render(w, http.StatusInternalServerError, Errorpage, nil)
				return
			}
		}
	}

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
		// app.render(w, http.StatusInternalServerError, Errorpage, nil)
		// return
	}

	err = app.SessionManager.RenewToken(w, r, userID)
	if err != nil {
		app.Logger.Error("renewtoken", "error", err)
		app.render(w, http.StatusInternalServerError, Errorpage, nil)
		return
	}

	http.Redirect(w, r, redirectUrl, http.StatusSeeOther)
}
