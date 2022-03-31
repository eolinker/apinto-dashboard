package apinto_dashboard

import "net/http"

type IModule interface {
	http.Handler
}

