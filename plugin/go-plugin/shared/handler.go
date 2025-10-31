package shared

import (
	"context"
	module "github.com/eolinker/apinto-dashboard/module"
	"github.com/eolinker/apinto-dashboard/pm3"
)

const PluginHandlerName = "apinto-dashboard"

type MiddlewareHandler interface {
	OnHttpRequest(r *module.MiddlewareRequest) module.MiddlewareResponseWriter
	onHttpResponse(r *module.MiddlewareRequest) module.MiddlewareResponseWriter
	Check(api pm3.ApiInfo) bool
}

type Module interface {
	Frontend() []pm3.FrontendAsset
	Apis() []pm3.Api
	Middleware() []Middleware
}
type MiddlewareHandFunc func(ctx context.Context, request *module.MiddlewareRequest, writer module.MiddlewareResponseWriter)

type Middleware interface {
	RequestHandler() (MiddlewareHandFunc, bool)
	ResponseHandler() (MiddlewareHandFunc, bool)
	Check(api pm3.ApiInfo) bool
}
