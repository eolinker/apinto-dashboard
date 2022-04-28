package apinto_dashboard

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/internal/template"
	"github.com/go-basic/uuid"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	SessionName = "SESSIONID"
	CallBack    = "callback"
)

type AccountHandler struct {
	userDetailsService IUserDetailsService
	sessionManager     *SessionManager
	serHandler         http.Handler
	blackList          map[string]bool
}

func (h *AccountHandler) logoutHandler(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie(SessionName)
	if err == nil {
		if sessionCookie.Value != "" {

			h.sessionManager.Delete(sessionCookie.Value)
			sessionCookie.Value = ""
			http.SetCookie(w, sessionCookie)
		}
	}
	h.toLogin(w, r, "/")

}

func (h *AccountHandler) toLogin(w http.ResponseWriter, request *http.Request, callback string) {
	vs := url.Values{}
	vs.Set(CallBack, callback)
	http.Redirect(w, request, fmt.Sprintf("/login?%s", vs.Encode()), http.StatusSeeOther)
}
func (h *AccountHandler) loginHandler(w http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		{
			template.Execute(w, "login", map[string]interface{}{
				"code":   0,
				CallBack: request.URL.Query().Get(CallBack),
			})
			return
		}
	case http.MethodPost:
		h.Post(w, request)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}
func (h *AccountHandler) Post(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	callback := r.FormValue(CallBack)
	if callback == "" {
		callback = "/"
	}
	userName := r.FormValue("username")
	if userName == "" {
		template.Execute(w, "login", map[string]interface{}{
			"code":    1,
			"message": "user name mush input",
		})
		return
	}
	password := r.FormValue("password")
	if userName == "" {
		template.Execute(w, "login", map[string]interface{}{
			"code":    2,
			"message": "password mash input",
		})
		return
	}

	userDetails, err := h.userDetailsService.LoadUserByUsername(userName)
	if err != nil {
		template.Execute(w, "login", map[string]interface{}{
			"code":    3,
			"message": "username not exist",
		})
		AddActivityLog("unknown", "login", "system", "登陆失败，用户不存在", []*Arg{{Key: "username", Value: userName}, {Key: "callback", Value: callback}})

		return
	}
	if !strings.EqualFold(password, userDetails.GetPassword()) {
		template.Execute(w, "login", map[string]interface{}{
			"code":    4,
			"message": "username or password wrong",
		})
		AddActivityLog("unknown", "login", "system", "登陆失败，密码错误", []*Arg{{Key: "username", Value: userName}, {Key: "callback", Value: callback}})

		return
	}
	if !userDetails.IsAccountNonExpired() {
		template.Execute(w, "login", map[string]interface{}{
			"code":    5,
			"message": "user is expired",
		})
		AddActivityLog(userDetails.GetUsername(), "login", "system", "登陆失败，用户已过期", []*Arg{{Key: "username", Value: userName}, {Key: "callback", Value: callback}})

		return
	}
	if !userDetails.IsEnabled() {
		template.Execute(w, "login", map[string]interface{}{
			"code":    6,
			"message": "user is not enabled",
		})
		AddActivityLog(userDetails.GetUsername(), "login", "system", "登陆失败，用户未激活", nil)

		return
	}
	if !userDetails.IsAccountNonLocked() {
		template.Execute(w, "login", map[string]interface{}{
			"code":    7,
			"message": "user is locked",
		})
		AddActivityLog(userDetails.GetUsername(), "login", "system", "登陆失败，用户已锁定", []*Arg{{Key: "username", Value: userName}, {Key: "callback", Value: callback}})

		return
	}
	if !userDetails.IsCredentialsNonExpired() {
		template.Execute(w, "login", map[string]interface{}{
			"code":    8,
			"message": "user is credentials expires",
		})
		AddActivityLog(userDetails.GetUsername(), "login", "system", "登陆失败，证书已过期", []*Arg{{Key: "username", Value: userName}, {Key: "callback", Value: callback}})

		return
	}
	sessionId := uuid.New()
	h.sessionManager.Set(sessionId, userDetails)
	cookie := &http.Cookie{
		Name:  SessionName,
		Value: sessionId,
		Path:  "/",
	}
	isRemember := r.FormValue("remember") == "remember-me"
	if isRemember {
		cookie.Expires = time.Now().Add(time.Hour * 24)
	}
	http.SetCookie(w, cookie)

	AddActivityLog(userDetails.GetUsername(), "login", "system", "成功登陆", []*Arg{{Key: "username", Value: userName}, {Key: "callback", Value: callback}})
	http.Redirect(w, r, callback, http.StatusFound)
	//h.serHandler.ServeHTTP(w, setUserDetailsToRequest(r, userDetails))
}
func (h *AccountHandler) Api(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie(SessionName)
	if err != nil {

		return
	}
	userDetails, has := h.sessionManager.Get(sessionCookie.Value)
	if !has {

		return
	}
	h.serHandler.ServeHTTP(w, setUserDetailsToRequest(r, userDetails))

}
func NewAccountHandler(userDetailsService IUserDetailsService, ser http.Handler, blacklist []string) http.Handler {

	srv := &http.ServeMux{}

	accountHandler := &AccountHandler{
		userDetailsService: userDetailsService,
		sessionManager:     NewSessionManager(),
		serHandler:         ser,
		blackList:          make(map[string]bool),
	}

	for _, p := range blacklist {
		srv.Handle(p, ser)
	}
	srv.Handle("/favicon.ico", ser)
	srv.HandleFunc("/login", accountHandler.loginHandler)
	srv.HandleFunc("/logout", accountHandler.logoutHandler)
	srv.HandleFunc("/api/", accountHandler.Api)
	srv.HandleFunc("/profession/", accountHandler.Api)
	srv.HandleFunc("/", accountHandler.View)
	return srv
}

func (h *AccountHandler) View(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie(SessionName)
	if err != nil {
		h.toLogin(w, r, r.RequestURI)
		return
	}
	userDetails, has := h.sessionManager.Get(sessionCookie.Value)
	if !has {
		h.toLogin(w, r, r.RequestURI)
		return
	}

	if !sessionCookie.Expires.IsZero() {
		sessionCookie.Expires = time.Now().Add(time.Hour * 24)
		http.SetCookie(w, sessionCookie) // 更新sesion过期时间
	}

	h.serHandler.ServeHTTP(w, setUserDetailsToRequest(r, userDetails))
}

func (h *AccountHandler) Clear() {

}
