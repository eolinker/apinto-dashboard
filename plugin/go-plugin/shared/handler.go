package shared

import (
	module "github.com/eolinker/apinto-dashboard/module"
	"net/http"
)

const PluginHandlerName = "apinto-dashboard"

type HttpHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type MiddlewareHandler interface {
	OnHttpRequest(r *module.MiddlewareRequest) module.MiddlewareResponseWriter
	onHttpResponse(r *module.MiddlewareRequest) module.MiddlewareResponseWriter
}
