package handler

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"sync"
	"time"

	"forum/internal/entities"
	"forum/internal/session"
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
		// Используем reflect для получения имени функции
		h = c.middlewares[i](h)
	}
	return h
}

func (c *Chain) ThenFunc(hf http.HandlerFunc) http.Handler {
	return c.Then(hf)
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

func (app *Application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.Logger.Info("Received request",
			"remote_addr", r.RemoteAddr,
			"protocol", r.Proto,
			"method", r.Method,
			"url", r.URL.RequestURI(),
		)

		next.ServeHTTP(w, r)
	})
}

func (app *Application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				err = fmt.Errorf("%s", err)
				app.Logger.Error("recover panic",
					"error", err,
					"stack_trace", string(debug.Stack()), // Включение стека
				)

				app.render(w, http.StatusInternalServerError, Errorpage, nil)
				return
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *Application) verifyCSRF(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем (все) запросы которые могут изменить данные
		if r.Method == http.MethodPost {
			sess, err := app.SessionManager.SessionStart(w, r)
			if err != nil {
				app.Logger.Error("verifyCSRF", "error", err)
				app.render(w, http.StatusInternalServerError, Errorpage, nil)
				return
			}
			sessionToken, ok := sess.Get(CsrfTokenSessionKey).(string)
			if !ok || sessionToken == "" {
				sessionToken = app.generateCSRFToken()
				if err := sess.Set(CsrfTokenSessionKey, sessionToken); err != nil {
					app.Logger.Error("get csrftoken", "error", err)
					app.render(w, http.StatusInternalServerError, Errorpage, nil)
					return
				}
			}
			requestToken := r.FormValue(CsrfTokenSessionKey)

			// Вставляем сессию в контекст запроса, чтобы другие хэндлеры могли её использовать
			ctx := context.WithValue(r.Context(), sessionContextKey, sess)
			r = r.WithContext(ctx)

			app.Logger.Debug("tokens in verifyCSRF", "request", requestToken, "session", sessionToken)
			if requestToken != sessionToken {
				app.Logger.Error("Invalid CSRF token", "error", err)
				app.render(w, http.StatusForbidden, Errorpage, nil)
				return
			}
		}

		// Продолжить выполнение запроса
		next.ServeHTTP(w, r)
	})
}

func (app *Application) sessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Извлекаем сессию из контекста
		sess, ok := r.Context().Value(sessionContextKey).(session.Session)
		if !ok {
			var err error
			sess, err = app.SessionManager.SessionStart(w, r)
			if err != nil {
				app.Logger.Error("SessionStart", "error", err)
				app.render(w, http.StatusInternalServerError, Errorpage, nil)
				return
			}
			app.Logger.Debug("session in sessionMiddleware", "session", sess)
		}

		// Если токен уже существует в сессии, не перезаписываем его
		token, ok := sess.Get(CsrfTokenSessionKey).(string)
		if !ok || token == "" {
			token = app.generateCSRFToken()
			if err := sess.Set(CsrfTokenSessionKey, token); err != nil {
				app.Logger.Error("get csrftoken from session", "error", err)
				app.render(w, http.StatusInternalServerError, Errorpage, nil)
				return
			}
		}

		// Вставляем токен и сессию в контекст запроса, чтобы другие хэндлеры могли её использовать
		ctx := context.WithValue(r.Context(), csrfTokenContextKey, token)
		ctx = context.WithValue(ctx, sessionContextKey, sess)
		r = r.WithContext(ctx)

		// Продолжить выполнение запроса
		next.ServeHTTP(w, r)
	})
}

func (app *Application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess := app.SessionFromContext(r)
		id, ok := sess.Get(AuthUserIDSessionKey).(int)
		if !ok || id == 0 {
			next.ServeHTTP(w, r)
			return
		}

		exists, err := app.Service.User.UserExists(id)
		if err != nil {
			app.Logger.Error("user exists", "error", err)
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
			return
		}

		if exists {
			ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

func (app *Application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			// Add the path that the user is trying to access to their session
			// data.
			sess := app.SessionFromContext(r)
			sess.Set(RedirectPathAfterLoginSessionKey, r.URL.Path)
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		// set the "Cache-Control: no-store" header so that pages
		// require authentication are not stored in the users browser cache (or
		// other intermediary cache).
		w.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}

func (app *Application) requireModeration(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess := app.SessionFromContext(r)
		userRole, ok := sess.Get(UserRoleSessionKey).(string)
		if !ok {
			app.Logger.Error("cannot extract user role from session in requireModeration")
			app.render(w, http.StatusInternalServerError, Errorpage, nil)
			return
		}
		if !(userRole == entities.RoleModerator || userRole == entities.RoleAdmin) {
			app.Logger.Warn("user role not moderator or admin in requireModeration")
			app.render(w, http.StatusForbidden, Errorpage, nil)
			return
		}

		next.ServeHTTP(w, r)
	})
}

type rateLimiter struct {
	visitors      sync.Map
	rate          int           // Количество запросов
	burstInterval time.Duration // Интервал времени для сброса
}

type visitor struct {
	lastSeen time.Time
	count    int
}

// NewRateLimiter создает новый rateLimiter
func NewRateLimiter(rate int, interval time.Duration) *rateLimiter {
	rl := &rateLimiter{
		rate:          rate,
		burstInterval: interval,
	}
	go rl.cleanupVisitors()
	return rl
}

// Middleware реализует ограничение скорости
func (app *Application) rateLimiting(next http.Handler) http.Handler {
	// Создаем Rate Limiter
	rl := NewRateLimiter(120, 1*time.Minute) // 60 запросов в минуту

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		v, _ := rl.visitors.LoadOrStore(ip, &visitor{lastSeen: time.Now(), count: 0})

		vis := v.(*visitor)
		vis.lastSeen = time.Now()

		if vis.count >= rl.rate {
			app.Logger.Error("Too many requests", "from", ip)
			app.render(w, http.StatusTooManyRequests, Errorpage, nil)
			return
		}

		vis.count++
		next.ServeHTTP(w, r)
	})
}

// cleanupVisitors очищает старые записи
func (rl *rateLimiter) cleanupVisitors() {
	for range time.NewTicker(rl.burstInterval).C {
		now := time.Now()
		rl.visitors.Range(func(key, value any) bool {
			vis := value.(*visitor)
			if now.Sub(vis.lastSeen) > rl.burstInterval {
				rl.visitors.Delete(key)
			}
			return true
		})
	}
}

// Middleware для добавления заголовков кэширования
func cacheControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Устанавливаем заголовки для кэширования
		w.Header().Set("Cache-Control", "public, max-age=10800") // Кэширование на 3 часа
		next.ServeHTTP(w, r)
	})
}
