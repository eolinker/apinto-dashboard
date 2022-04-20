package apinto_dashboard

import (
	"net/http"
)

type AccountHandler struct {
	server             http.Handler
	userDetailsService IUserDetailsService
	sessionManager     SessionManager
}

func NewAccountHandler(userDetailsService IUserDetailsService) http.Handler {
	return &AccountHandler{userDetailsService: userDetailsService}
}

func (h *AccountHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
