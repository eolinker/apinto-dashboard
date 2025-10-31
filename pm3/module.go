package pm3

import (
	"github.com/gin-gonic/gin"
)

type FrontendPluginItem struct {
}

type FrontendAsset struct {
	Path        string
	HandlerFunc gin.HandlerFunc
}
type ApiInfo struct {
	Authority ApiAuthority
	Access    string
	Method    string
	Path      string
}
type Api struct {
	Authority   ApiAuthority
	Access      string
	Method      string
	Path        string
	HandlerFunc gin.HandlerFunc
}

func (A Api) Info() ApiInfo {
	return ApiInfo{
		Authority: A.Authority,
		Access:    A.Access,
		Method:    A.Method,
		Path:      A.Path,
	}
}

type MiddlewareChecker interface {
	Check(api ApiInfo) bool
}
type Middleware interface {
	Handle(ginCtx *gin.Context)
	MiddlewareChecker
}

type Module interface {
	Name() string
	Frontend() []FrontendAsset
	Apis() []Api
	Middleware() []Middleware
	Support() (ProviderSupport, bool)
}
