package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// Middleware type for handling HTTP requests
type Middleware func(http.Handler) http.Handler

// Chain struct to handle multiple middleware functions
type Chain struct {
	middlewares []Middleware
}

// New creates a new chain with given middlewares
func New(middlewares ...Middleware) *Chain {
	return &Chain{middlewares: middlewares}
}

// Append adds more middlewares to the chain
func (c *Chain) Append(middlewares ...Middleware) *Chain {
	return &Chain{middlewares: append(c.middlewares, middlewares...)}
}

// Then applies all middleware to the given handler
func (c *Chain) Then(h http.Handler) http.Handler {
	for i := len(c.middlewares) - 1; i >= 0; i-- {
		// for i := 0; i < len(c.middlewares); i++ {
		// Используем reflect для получения имени функции
		h = c.middlewares[i](h)
	}
	return h
}

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) verifyCSRF(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем все запросы которые могут изменить данные
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			sess := app.sessionManager.SessionStart(w, r)
			sessionToken := sess.Get("token").(string)
			err := sess.Delete("token")
			if err != nil {
				app.serverError(w, err)
			}
			requestToken := r.FormValue("token")

			if sessionToken == "" || requestToken != sessionToken {
				http.Error(w, "Invalid CSRF token", http.StatusForbidden)
				return
			}
		}

		// Продолжить выполнение запроса
		next.ServeHTTP(w, r)
	})
}

func (app *application) sessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess := app.sessionManager.SessionStart(w, r)

		// Генерация нового CSRF-токена
		token := app.generateCSRFToken()
		sess.Set("token", token)

		// Вставка токена в контекст запроса
		ctx := context.WithValue(r.Context(), "csrfToken", token)
		r = r.WithContext(ctx)

		user := sess.Get("username")
		role := sess.Get("role")

		if user == nil && requiresAuth(r.URL.Path) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		// Проверка ролей для защищённых маршрутов
		if role != nil && hasAccess(role, r.URL.Path) {
			// http.Error(w, "Forbidden", http.StatusForbidden)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Продолжить выполнение запроса
		next.ServeHTTP(w, r)
	})
}

func requiresAuth(path string) bool {
	// Определить, требует ли страница авторизации
	return strings.HasPrefix(path, "/post/create") || strings.HasPrefix(path, "/user/logout") || strings.HasPrefix(path, "/moderation")
}

func hasAccess(role interface{}, path string) bool {
	// Определяем доступность страниц для различных ролей
	switch role {
	case "admin":
		return true // Админы имеют доступ ко всему
	case "moderator":
		// Модеры имеют доступ к страницам модерации
		// if strings.HasPrefix(path, "/moderation") {
		return true
		// }
	case "user":
		// Обычные пользователи могут создавать посты, но не имеют доступа к модерации
		if strings.HasPrefix(path, "/moderation") || strings.HasPrefix(path, "/user/signup") || strings.HasPrefix(path, "/user/login") {
			return false
		}
	// Гостям доступ ограничен, они могут только просматривать
	default:
		return false
	}
	return false
}
