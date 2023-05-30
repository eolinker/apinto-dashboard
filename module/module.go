package apinto_module

import (
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	RouterTypeModule    = "module"
	RouterTypeApi       = "api"
	RouterTypeAsset     = "assets"
	RouterTypeOpenApi   = "openapi"
	RouterTypeAnonymous = "anonymous"
	RouterTypeSensitive = "sensitive"
)

var (
	RouterLabelAssets    = []string{RouterTypeAsset}
	RouterLabelModule    = []string{RouterTypeModule}
	RouterLabelApi       = []string{RouterTypeApi}
	RouterLabelOpenApi   = []string{RouterTypeOpenApi}
	RouterLabelAnonymous = []string{RouterTypeAnonymous}
)

type Module interface {
	Name() string
	Routers() (Routers, bool)         // 是否提供路由
	Middleware() (Middleware, bool)   //是否提供拦截器
	Support() (ProviderSupport, bool) //是否提供provider
}
type CoreModule interface {
	Module
}
type HandlerFunc = gin.HandlerFunc

type RouterInfo struct {
	Method      string
	Path        string
	Handler     string
	Labels      []string
	HandlerFunc []HandlerFunc
	Replaceable bool
	labels      map[string]struct{}
}

var (
	hashSetValue = struct{}{}
)

func defaultType(path string) string {
	if strings.HasPrefix(path, "/api") {
		return RouterTypeApi
	}
	if strings.HasPrefix(path, "/api1") {
		return RouterTypeOpenApi
	}
	return RouterTypeModule
}

func (r RouterInfo) getLabels() map[string]struct{} {
	if r.labels != nil {
		return r.labels
	}
	lbs := make(map[string]struct{})
	lbs[r.Path] = hashSetValue
	lbs[r.Method] = hashSetValue
	if len(r.Labels) == 0 {
		lbs[defaultType(r.Path)] = struct{}{}
	} else {
		for _, l := range r.Labels {
			lbs[l] = hashSetValue
		}
	}

	r.labels = lbs
	return r.labels
}

type RoutersInfo []RouterInfo

func (rs RoutersInfo) Find(method, path string) (*RouterInfo, bool) {
	for _, r := range rs {
		if r.Method == method && r.Path == path {
			return &r, true
		}
	}
	return nil, false
}

type MiddlewareHandler struct {
	Name        string
	Rule        MiddlewareChecker
	Handler     HandlerFunc
	Replaceable bool
}

func (m *MiddlewareHandler) checkRouter(router RouterInfo) bool {
	return m.Rule.checkRouter(router)
}

type MiddlewareChecker interface {
	checkRouter(router RouterInfo) bool
}
type MiddlewareRule []string

func (rv MiddlewareRule) checkRouter(router RouterInfo) bool {
	lbs := router.getLabels()
	for _, v := range rv {
		if _, has := lbs[v]; !has {
			return false
		}
	}
	return true
}

type MiddlewareRules []MiddlewareRule

func CreateMiddlewareRules(rs interface{}) MiddlewareChecker {

	switch vs := rs.(type) {
	case []string:
		return MiddlewareRule(vs)
	case [][]string:
		rl := make(MiddlewareRules, 0, len(vs))
		for _, r := range vs {
			rl = append(rl, r)
		}
		return rl
	case string:
		return MiddlewareRule{vs}
	default:
		return MiddlewareRules{}
	}

}
func (mrs MiddlewareRules) checkRouter(router RouterInfo) bool {
	for _, mr := range mrs {
		if mr.checkRouter(router) {
			return true
		}
	}
	return false
}

type Routers interface {
	RoutersInfo() RoutersInfo
}
type Middleware interface {
	MiddlewaresInfo() []MiddlewareHandler
}
type Cargo struct {
	Value string
	Title string
}

func (c *Cargo) Export() *CargoItem {
	return &CargoItem{
		Value: c.Value,
		Title: c.Title,
	}
}

type CargoStatus int

const (
	None    = iota // 无效
	Offline        // 未上线
	Online         // 已上线
)

type Provider interface {
	Provide(namespaceId int) []Cargo
}
type ProviderStatus interface {
	Status(key string, namespaceId int, cluster string) (CargoStatus, string)
}
type ProviderSupport interface {
	Provider() map[string]Provider
	ProviderStatus
}
