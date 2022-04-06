package apinto_dashboard

import "errors"

var (
	ErrorNotLogin = errors.New("not login")
	ErrorUsernameNotFound = errors.New("username not found")
	ErrorUserDetailsServiceNeed = errors.New("need UserDetailsService")
)