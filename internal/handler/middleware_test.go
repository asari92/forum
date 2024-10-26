package handler

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"forum/pkg/assert"
)

func TestSecureHeaders(t *testing.T) {
	// Initialize a new httptest.ResponseRecorder and dummy http.Request.
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock HTTP handler that we can pass to our secureHeaders
	// middleware, which writes a 200 status code and an "OK" response body.
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Pass the mock HTTP handler to our secureHeaders middleware. Because
	// secureHeaders *returns* a http.Handler we can call its ServeHTTP()
	// method, passing in the http.ResponseRecorder and dummy http.Request to
	// execute it.
	secureHeaders(next).ServeHTTP(rr, r)

	// Call the Result() method on the http.ResponseRecorder to get the results
	// of the test.
	rs := rr.Result()

	// Check that the middleware has correctly set the Content-Security-Policy
	// header on the response.
	expectedValue := "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com"
	assert.Equal(t, rs.Header.Get("Content-Security-Policy"), expectedValue)

	// Check that the middleware has correctly set the Referrer-Policy
	// header on the response.
	expectedValue = "origin-when-cross-origin"
	assert.Equal(t, rs.Header.Get("Referrer-Policy"), expectedValue)

	// Check that the middleware has correctly set the X-Content-Type-Options
	// header on the response.
	expectedValue = "nosniff"
	assert.Equal(t, rs.Header.Get("X-Content-Type-Options"), expectedValue)

	// Check that the middleware has correctly set the X-Frame-Options header
	// on the response.
	expectedValue = "deny"
	assert.Equal(t, rs.Header.Get("X-Frame-Options"), expectedValue)

	// Check that the middleware has correctly set the X-XSS-Protection header
	// on the response
	expectedValue = "0"
	assert.Equal(t, rs.Header.Get("X-XSS-Protection"), expectedValue)

	// Check that the middleware has correctly called the next handler in line
	// and the response status code and body are as expected.
	assert.Equal(t, rs.StatusCode, http.StatusOK)

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")
}

func TestVerifyCSRF(t *testing.T) {
	app := newTestApplication(t)

	// Создаем тестовый POST-запрос с CSRF-токеном
	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodPost, "/post/create", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Устанавливаем фальшивую сессию с токеном
	sess := app.SessionManager.SessionStart(rr, r)
	token := app.generateCSRFToken()
	err = sess.Set(CsrfTokenSessionKey, token)
	if err != nil {
		t.Fatal(err)
	}

	// Теперь нужно вытащить куку из ответа rr и установить её в следующий запрос
	cookie := rr.Result().Cookies()
	for _, c := range cookie {
		r.AddCookie(c) // Передаем куки в запрос
	}

	r.PostForm = map[string][]string{
		CsrfTokenSessionKey: {token},
	}

	// Запускаем middleware и проверяем статус
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	app.verifyCSRF(next).ServeHTTP(rr, r)
	// в миддлеверке первым делом создается сессия
	rs := rr.Result()
	if rs.StatusCode != http.StatusOK {
		t.Errorf("expected status %d; got %d", http.StatusOK, rs.StatusCode)
	}
}

func TestSessionMiddleware(t *testing.T) {
	app := newTestApplication(t)

	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess := app.SessionFromContext(r)
		if sess == nil {
			t.Fatal("expected a session in context")
		}

		token := r.Context().Value(csrfTokenContextKey)
		if token == "" {
			t.Fatal("expected CSRF token in context")
		}
		w.WriteHeader(http.StatusOK)
	})

	app.sessionMiddleware(next).ServeHTTP(rr, r)

	rs := rr.Result()
	if rs.StatusCode != http.StatusOK {
		t.Errorf("expected status %d; got %d", http.StatusOK, rs.StatusCode)
	}
}

func TestRequireAuthentication200(t *testing.T) {
	app := newTestApplication(t)

	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/post/create", nil)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
	r = r.WithContext(ctx)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	app.requireAuthentication(next).ServeHTTP(rr, r)

	rs := rr.Result()
	if rs.StatusCode != http.StatusOK {
		t.Errorf("expected status %d; got %d", http.StatusOK, rs.StatusCode)
	}
}

func TestRequireAuthentication303(t *testing.T) {
	app := newTestApplication(t)

	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/post/create", nil)
	if err != nil {
		t.Fatal(err)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	sess := app.SessionManager.SessionStart(rr, r)
	ctx := context.WithValue(r.Context(), sessionContextKey, sess)
	r = r.WithContext(ctx)

	app.requireAuthentication(next).ServeHTTP(rr, r)

	rs := rr.Result()
	if rs.StatusCode != http.StatusSeeOther {
		t.Errorf("expected status %d; got %d", http.StatusSeeOther, rs.StatusCode)
	}
}

func TestAuthenticate(t *testing.T) {
	app := newTestApplication(t)

	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	sess := app.SessionManager.SessionStart(rr, r)
	err = sess.Set(AuthUserIDSessionKey, 1)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.WithValue(r.Context(), sessionContextKey, sess)
	r = r.WithContext(ctx)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAuthenticated := r.Context().Value(isAuthenticatedContextKey).(bool)
		if !isAuthenticated {
			t.Fatal("expected user to be authenticated")
		}
		w.WriteHeader(http.StatusOK)
	})

	app.authenticate(next).ServeHTTP(rr, r)

	rs := rr.Result()
	if rs.StatusCode != http.StatusOK {
		t.Errorf("expected status %d; got %d", http.StatusOK, rs.StatusCode)
	}
}
