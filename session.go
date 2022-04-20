package apinto_dashboard

import "sync"

type SessionManager struct {
	lock     sync.RWMutex
	users    map[string]string
	sessions map[string]UserDetails
}

func (sm *SessionManager) Get(session string) UserDetails {
	sm.lock.RLock()
	defer sm.lock.RUnlock()

}
