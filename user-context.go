package apinto_dashboard

import (
	"context"
	"net/http"
)

type userDetailsKey struct{}

// UserDetailsKey  is the request context key under which UserDetails are stored.
var UserDetailsKey = userDetailsKey{}

func UserDetailsFromRequest(req *http.Request) (UserDetails,error) {
	value,ok := req.Context().Value(UserDetailsKey).(UserDetails)
	if !ok{
		return nil, ErrorNotLogin
	}
	return value,nil
}
func setUserDetailsToRequest(req *http.Request, details UserDetails) *http.Request{
	if details == nil{
		return req
	}
	ctx := req.Context()
	ctx = context.WithValue(ctx, UserDetailsKey, details)
	req = req.WithContext(ctx)
	return req
}