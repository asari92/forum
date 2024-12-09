package handler

import (
	"context"
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

type userInfo struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (app *Application) oauthGoogleLogin(w http.ResponseWriter, r *http.Request) {
	oauthState := app.generateCSRFToken()

	sess := app.SessionFromContext(r)
	sess.Set(GoogleOAuthStateSessionKey, oauthState)
	app.Logger.Info("oauth state was generated successfull", "oauthState", oauthState)

	googleOauthConfig.ClientID = app.Config.GoogleClientID
	googleOauthConfig.ClientSecret = app.Config.GoogleClientSecret
	googleOauthConfig.RedirectURL = app.Config.GoogleClientCallbackURL
	url := googleOauthConfig.AuthCodeURL(oauthState, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// func generateStateOauthCookie(sess session.Session) string {
// 	b := make([]byte, 16)
// 	if _, err := rand.Read(b); err != nil {
// 		app.Logger.Error("Failed to generate random state", "error", err)
// 		return ""
// 	}
// 	state := base64.URLEncoding.EncodeToString(b)
// 	sess.Set(GoogleOAuthStateSessionKey, state)
// 	return state
// }

// Callback от Google
func (app *Application) oauthGoogleCallback(w http.ResponseWriter, r *http.Request) {
	sess := app.SessionFromContext(r)

	sessionState, ok := sess.Get(GoogleOAuthStateSessionKey).(string)
	if !ok || sessionState != r.FormValue("state") {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	code := r.FormValue("code")
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	userInfo, err := fetchGoogleUserInfo(token)
	if err != nil {
		http.Error(w, "Failed to fetch GoogleUserInfo", http.StatusInternalServerError)
		return
	}

	userID, err := app.Service.User.OauthAuthenticate(userInfo.Email)
	if err != nil {
		if errors.Is(err, entities.ErrInvalidCredentials) {
		} else {
			app.Logger.Error("get id Authenticate user", "error", err)
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
			return
		}
	}

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

	err = sess.Delete(GoogleOAuthStateSessionKey)
	if err != nil {
		app.Logger.Error("Session error during delete google oauth state", "error", err)
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

func fetchGoogleUserInfo(token *oauth2.Token) (*userInfo, error) {
	client := googleOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo userInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}
