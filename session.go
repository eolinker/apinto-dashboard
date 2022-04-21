package apinto_dashboard

import "sync"

type SessionManager struct {
	lock     sync.RWMutex
	users    map[string]string
	sessions map[string]UserDetails
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		lock:     sync.RWMutex{},
		users:    make(map[string]string),
		sessions: make(map[string]UserDetails),
	}
}

func (sm *SessionManager) Get(session string) (UserDetails, bool) {
	sm.lock.RLock()
	defer sm.lock.RUnlock()
	details, has := sm.sessions[session]
	return details, has
}
func (sm *SessionManager) Delete(session string) {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	old, has := sm.sessions[session]
	if has {
		delete(sm.sessions, session)
		if old != nil {
			delete(sm.users, old.GetUsername())
		}
	}
}
func (sm *SessionManager) Set(session string, details UserDetails) {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	old, has := sm.users[details.GetUsername()]
	if has {
		delete(sm.sessions, old)
	}
	sm.users[details.GetUsername()] = session
	sm.sessions[session] = details

}
