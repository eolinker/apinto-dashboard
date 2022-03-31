package security

type UserDetails struct {
	username string
	password string
	authorities map[string]interface{}
}

func NewUserDetails(username string, password string, authorities map[string]interface{}) *UserDetails {
	return &UserDetails{username: username, password: password, authorities: authorities}
}

func (u *UserDetails) GetUsername() string {
	return u.username
}

func (u *UserDetails) GetAuthorities() map[string]interface{} {
	return u.authorities
}

func (u *UserDetails) GetPassword() string {
	return u.password
}

func (u *UserDetails) IsAccountNonExpired() bool {
	return true
}

func (u *UserDetails) IsAccountNonLocked() bool {
	return true
}

func (u *UserDetails) IsCredentialsNonExpired() bool {
	return true
}

func (u *UserDetails) IsEnabled() bool {
	return true
}

