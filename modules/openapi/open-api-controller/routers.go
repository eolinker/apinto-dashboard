package open_api_controller

import (
	"github.com/eolinker/apinto-dashboard/pm3"
	"github.com/eolinker/eosc"
)

type IRouter interface {
	Name() string
	Apis() []pm3.Api
}

var defaultRouterManager = NewRouterManager()

type RouterManager struct {
	routers eosc.Untyped[string, IRouter]
}

func (r *RouterManager) Register(router IRouter) {
	r.routers.Set(router.Name(), router)
}

func (r *RouterManager) GetRouter(name string) (IRouter, bool) {
	return r.routers.Get(name)
}

func (r *RouterManager) AllRouters() []IRouter {
	return r.routers.List()
}

func NewRouterManager() *RouterManager {
	return &RouterManager{
		routers: eosc.BuildUntyped[string, IRouter](),
	}
}

func RegisterRouter(router IRouter) {
	defaultRouterManager.Register(router)
}

func GetRouter(name string) (IRouter, bool) {
	return defaultRouterManager.GetRouter(name)
}

func AllRouters() []IRouter {
	return defaultRouterManager.AllRouters()
}
