package apinto_dashboard


type UserDetails interface {

	// GetUsername  Returns the username used to authenticate the user.
	GetUsername()string
	// GetAuthorities Returns the authorities granted to the user.
	GetAuthorities()map[string]interface{}
	// GetPassword Returns the password used to authenticate the user.
	GetPassword()string
	// IsAccountNonExpired Indicates whether the user's account has expired.
	IsAccountNonExpired() bool
	// IsAccountNonLocked Indicates whether the user is locked or unlocked.
	IsAccountNonLocked() bool
	// IsCredentialsNonExpired Indicates whether the user's credentials (password) has expired.
	IsCredentialsNonExpired() bool
	// IsEnabled Indicates whether the user is enabled or disabled.
	IsEnabled() bool
}

type IUserDetailsService interface {
	// LoadUserByUsername Locates the user based on the username.
	LoadUserByUsername(username string) (UserDetails,error)
}