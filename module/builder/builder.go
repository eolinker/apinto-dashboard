package builder

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/module"
	"github.com/eolinker/apinto-dashboard/pm3"
	"github.com/eolinker/eosc/log"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

var (
	systemModules           []pm3.Module
	indexHtmlHandle         gin.HandlerFunc
	frontendAssetMiddleware []gin.HandlerFunc
)

func AddAssetMiddleware(hs ...gin.HandlerFunc) {
	frontendAssetMiddleware = append(frontendAssetMiddleware, hs...)
}
func SetIndexHtml(h gin.HandlerFunc) {
	indexHtmlHandle = h
}
func AddSystemModule(ms ...pm3.Module) {
	systemModules = append(systemModules, ms...)
}

type FrontendAssets struct {
	pm3.FrontendAsset
	plugin string
}
type Middleware struct {
	pm3.Middleware
	plugin string
}
type Api struct {
	pm3.Api
	plugin string
}
type ModuleBuilder struct {
	lock sync.Mutex

	engine *gin.Engine

	frontendAssets      []FrontendAssets
	middlewares         []Middleware
	apis                []Api
	optionHandlers      []apinto_module.IFilterOptionHandler
	providerBuilder     *apinto_module.ProviderBuilder
	filterOptionManager apinto_module.IFilterOptionHandlerManager

	noRouter *noRouterPrefix
}

func NewModuleBuilder(engine *gin.Engine, filterOptionManager ...apinto_module.IFilterOptionHandlerManager) *ModuleBuilder {

	routes := engine.Routes()
	if len(routes) > 0 {
		for _, r := range routes {
			log.DebugF("unknown %s %s %s", r.Handler, r.Method, r.Path)
		}
		panic("unknown")
	}
	b := &ModuleBuilder{
		engine:          engine,
		providerBuilder: apinto_module.NewProviderBuilder(),
		noRouter:        newNoRouterPrefix(),
	}
	engine.GET("/", indexHtmlHandle)
	engine.NoRoute(b.noRouter.noRouterHandleFunc)
	b.Append(systemModules...)

	if len(filterOptionManager) > 0 {
		b.filterOptionManager = filterOptionManager[0]
	} else {
		b.filterOptionManager = nil
	}

	return b
}

func (b *ModuleBuilder) Append(module ...pm3.Module) *ModuleBuilder {
	if len(module) > 0 {

		b.lock.Lock()
		defer b.lock.Unlock()

		for _, sm := range module {

			if hs, ok := sm.(apinto_module.IFilterOptionHandlerSupport); ok {
				b.optionHandlers = append(b.optionHandlers, hs.FilterOptionHandler()...)
			}
			for _, f := range sm.Frontend() {
				b.frontendAssets = append(b.frontendAssets, FrontendAssets{
					FrontendAsset: f,
					plugin:        sm.Name(),
				})
			}
			for _, f := range sm.Apis() {
				b.apis = append(b.apis, Api{
					Api:    f,
					plugin: sm.Name(),
				})
			}
			for _, f := range sm.Middleware() {
				b.middlewares = append(b.middlewares, Middleware{
					Middleware: f,
					plugin:     sm.Name(),
				})
			}

			if p, ok := sm.Support(); ok {
				b.providerBuilder.Append(sm.Name(), p)
			}
		}

	}
	return b
}

// Build 构建模块路由、providers等
func (b *ModuleBuilder) Build() (httpHandler http.Handler, iproviders apinto_module.IProviders, errResult error) {
	defer func() {
		if e := recover(); e != nil {
			errResult = fmt.Errorf("%v", e)
		}
	}()
	b.lock.Lock()
	defer b.lock.Unlock()

	assetGroup := b.engine.Group("/").Use(frontendAssetMiddleware...)
	for _, as := range b.frontendAssets {
		path := apinto_module.StaticRouter(as.Path)
		b.noRouter.addFrontendAssets(path)

		log.DebugF("add assets %s", path)
		assetGroup.GET(path, as.HandlerFunc)
		assetGroup.HEAD(path, as.HandlerFunc)
	}
	apigroup := b.engine.Group("/")
	for _, api := range b.apis {
		mids := make([]gin.HandlerFunc, 0, len(b.middlewares)+2)
		mids = append(mids, apinto_module.ModuleNameHandler(api.plugin), apinto_module.ApiInfoSetHandler(api.Api))
		for _, m := range b.middlewares {
			if m.Check(api.Info()) {
				mids = append(mids, m.Handle)
			}
		}
		b.noRouter.addApi(api.Path)

		mids = append(mids, api.HandlerFunc)
		apigroup.Handle(api.Method, api.Path, mids...)
	}
	b.resetFiltersHandler()

	return b.engine.Handler(), b.providerBuilder.Build(), nil
}

func (b *ModuleBuilder) resetFiltersHandler() {
	if b.filterOptionManager == nil {
		return
	}
	handlers := common.SliceToMap(b.optionHandlers, func(t apinto_module.IFilterOptionHandler) string {
		return t.Name()
	})

	b.filterOptionManager.ResetFilterOptionHandlers(handlers)

}
