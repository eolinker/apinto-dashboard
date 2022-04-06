package security

import (
	apinto_dashboard "github.com/eolinker/apinto-dashboard"
	"sync"
)

type UserDetailsService struct {
	users map[string]apinto_dashboard.UserDetails
	locker sync.RWMutex
}

func NewUserDetailsService() *UserDetailsService {
	return &UserDetailsService{users: make(map[string]apinto_dashboard.UserDetails)}
}

func (us *UserDetailsService) LoadUserByUsername(username string) (apinto_dashboard.UserDetails, error) {
	us.locker.RLock()
	defer us.locker.RUnlock()
	user,has:= us.users[username]
	if has{
		return user,nil
	}
	return nil,apinto_dashboard.ErrorUsernameNotFound
}
func (us *UserDetailsService) Add(userDetails apinto_dashboard.UserDetails)  {
	if userDetails == nil{
		return
	}
	us.locker.Lock()
	defer us.locker.Unlock()
	us.users[userDetails.GetUsername()] = userDetails
}