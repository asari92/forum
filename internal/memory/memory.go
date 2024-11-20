package memory

import (
	"container/list"
	"fmt"
	"sync"
	"time"

	"forum/internal/session"
)

var pder = &Provider{list: list.New(), userSessions: map[int]string{}}

func init() {
	pder.sessions = make(map[string]*list.Element, 0)
	session.Register("memory", pder)
}

type SessionStore struct {
	sid          string                      // unique session id
	timeAccessed time.Time                   // last access time
	value        map[interface{}]interface{} // session value stored inside
}

type Provider struct {
	lock         sync.Mutex               // lock
	sessions     map[string]*list.Element // save in memory
	userSessions map[int]string           // userID -> sessionID
	list         *list.List               // gc
}

func (st *SessionStore) Set(key, value interface{}) error {
	err := pder.SessionUpdate(st.sid)
	if err != nil {
		return fmt.Errorf("failed to set the session value: %v", err)
	}
	st.value[key] = value
	return nil
}

func (st *SessionStore) Get(key interface{}) interface{} {
	err := pder.SessionUpdate(st.sid)
	if err != nil {
		fmt.Printf("update error when getting value from session: %v", err)
	}

	if v, ok := st.value[key]; ok {
		return v
	} else {
		return nil
	}
}

func (st *SessionStore) GetAll() map[interface{}]interface{} {
	err := pder.SessionUpdate(st.sid)
	if err != nil {
		fmt.Printf("update error when getting all values from session: %v", err)
	}
	return st.value
}

func (st *SessionStore) Delete(key interface{}) error {
	err := pder.SessionUpdate(st.sid)
	if err != nil {
		return fmt.Errorf("failed to delete the session value: %v", err)
	}
	delete(st.value, key)
	return nil
}

func (st *SessionStore) SessionID() string {
	return st.sid
}

func (pder *Provider) SessionInit(sid string) (session.Session, error) {
	pder.lock.Lock()
	defer pder.lock.Unlock()
	v := make(map[interface{}]interface{}, 0)
	newsess := &SessionStore{sid: sid, timeAccessed: time.Now(), value: v}
	element := pder.list.PushBack(newsess)
	pder.sessions[sid] = element
	return newsess, nil
}

func (pder *Provider) SessionRead(sid string) (session.Session, error) {
	if element, ok := pder.sessions[sid]; ok {
		return element.Value.(*SessionStore), nil
	} else {
		sess, err := pder.SessionInit(sid)
		return sess, err
	}
}

func (pder *Provider) SessionDestroy(sid string) error {
	if element, ok := pder.sessions[sid]; ok {
		delete(pder.sessions, sid)
		pder.list.Remove(element)
		return nil
	}
	return nil
}

func (pder *Provider) SessionGC(maxlifetime int64) {
	pder.lock.Lock()
	defer pder.lock.Unlock()

	for {
		element := pder.list.Back()
		if element == nil {
			break
		}
		if (element.Value.(*SessionStore).timeAccessed.Unix() + maxlifetime) < time.Now().Unix() {
			pder.list.Remove(element)
			delete(pder.sessions, element.Value.(*SessionStore).sid)
		} else {
			break
		}
	}
}

func (pder *Provider) SessionUpdate(sid string) error {
	pder.lock.Lock()
	defer pder.lock.Unlock()

	if element, ok := pder.sessions[sid]; ok {
		element.Value.(*SessionStore).timeAccessed = time.Now()
		pder.list.MoveToFront(element)
		return nil
	}

	return fmt.Errorf("session not found")
}

func (pder *Provider) SessionBindUser(userID int, sid string) error {
	pder.lock.Lock()
	defer pder.lock.Unlock()

	// Уничтожаем старую сессию, если она есть
	if oldSid, ok := pder.userSessions[userID]; ok {
		if err := pder.SessionDestroy(oldSid); err != nil {
			return err
		}
	}
	// Связываем новую сессию
	pder.userSessions[userID] = sid
	return nil
}
