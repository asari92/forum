package session

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gofrs/uuid"
)

type Manager struct {
	cookieName  string     // private cookiename
	lock        sync.Mutex // protects session
	provider    Provider
	maxlifetime int64
}

type Provider interface {
	SessionInit(sid string) (Session, error)
	SessionRead(sid string) (Session, error)
	SessionDestroy(sid string) error
	SessionGC(maxLifeTime int64)
	SessionBindUser(userID int, sessionID string) error
}

type Session interface {
	Set(key, value interface{}) error    // set session value
	Get(key interface{}) interface{}     // get session value
	GetAll() map[interface{}]interface{} // get all session values
	Delete(key interface{}) error        // delete session value
	SessionID() string                   // back current sessionID
}

var provides = make(map[string]Provider)

func NewManager(provideName, cookieName string, maxlifetime int64) (*Manager, error) {
	provider, ok := provides[provideName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provide %q (forgotten import?)", provideName)
	}
	return &Manager{provider: provider, cookieName: cookieName, maxlifetime: maxlifetime}, nil
}

func (manager *Manager) sessionId() string {
	id, err := uuid.NewV4()
	if err != nil {
		return ""
	}
	return id.String()
}

func (manager *Manager) SessionStart(w http.ResponseWriter, r *http.Request) (Session, error) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		sid := manager.sessionId()
		session, err := manager.provider.SessionInit(sid)
		if err != nil {
			return nil, err
		}
		// Установка cookie с флагами HttpOnly и Secure
		cookie := http.Cookie{
			Name:     manager.cookieName,
			Value:    url.QueryEscape(sid),
			Path:     "/",
			HttpOnly: true, // HttpOnly защищает от XSS-атак
			Secure:   true, // Secure защищает от передачи через HTTP (только HTTPS)
			MaxAge:   int(manager.maxlifetime),
			SameSite: http.SameSiteLaxMode,
		}
		http.SetCookie(w, &cookie)
		return session, nil
	} else {
		sid, _ := url.QueryUnescape(cookie.Value)
		session, _ := manager.provider.SessionRead(sid)
		return session, nil
	}
}

func (manager *Manager) RenewToken(w http.ResponseWriter, r *http.Request, userID int) error {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	// Получаем текущий сессионный ID из куки
	oldCookie, err := r.Cookie(manager.cookieName)
	if err != nil || oldCookie.Value == "" {
		return fmt.Errorf("session: no existing session")
	}

	// Читаем текущую сессию
	oldSid, _ := url.QueryUnescape(oldCookie.Value)
	oldSession, err := manager.provider.SessionRead(oldSid)
	if err != nil {
		return err
	}

	// Генерируем новый сессионный ID
	newSid := manager.sessionId()

	// Инициализируем новую сессию с новым сессионным ID
	newSession, err := manager.provider.SessionInit(newSid)
	if err != nil {
		return err
	}

	// Копируем данные из старой сессии в новую
	for key, value := range oldSession.GetAll() {
		newSession.Set(key, value)
	}

	// Уничтожаем старую сессию
	err = manager.provider.SessionDestroy(oldSid)
	if err != nil {
		return err
	}

	// Привязываем новую сессию к текущему пользователю
	err = manager.provider.SessionBindUser(userID, newSid)
	if err != nil {
		return fmt.Errorf("session bind user: %w", err)
	}

	// Устанавливаем новую куку с новым сессионным ID
	newCookie := http.Cookie{
		Name:     manager.cookieName,
		Value:    url.QueryEscape(newSid),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   int(manager.maxlifetime),
	}
	http.SetCookie(w, &newCookie)

	return nil
}

// Register makes a session provider available by the provided name.
// If a Register is called twice with the same name or if the driver is nil,
// it panics.
func Register(name string, provider Provider) {
	if provider == nil {
		panic("session: Register provider is nil")
	}
	if _, dup := provides[name]; dup {
		panic("session: Register called twice for provider " + name)
	}
	provides[name] = provider
}

// Destroy sessionid
func (manager *Manager) SessionDestroy(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		return
	} else {
		manager.lock.Lock()
		defer manager.lock.Unlock()
		manager.provider.SessionDestroy(cookie.Value)
		expiration := time.Now()
		cookie := http.Cookie{Name: manager.cookieName, Path: "/", HttpOnly: true, Expires: expiration, MaxAge: -1}
		http.SetCookie(w, &cookie)
	}
}

func (manager *Manager) GC() {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.provider.SessionGC(manager.maxlifetime)
	time.AfterFunc(time.Duration(manager.maxlifetime), func() { manager.GC() })
}
