package apinto_module

import (
	"errors"
	"fmt"
	"github.com/eolinker/eosc/log"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

var (
	systemMiddlewares []MiddlewareHandler
)

func AddSystemMiddleware(handler MiddlewareHandler) {
	systemMiddlewares = append(systemMiddlewares, handler)
}

type ModuleBuilder struct {
	lock              sync.Mutex
	modules           map[string]Module
	middlewaresAppend []MiddlewareHandler
	coreMiddleware    []MiddlewareHandler

	providerBuilder    *ProviderBuilder
	engine             *gin.Engine
	core               CoreModule
	alternativeRouters map[string]RouterInfo

	middlewares []MiddlewareHandler

	filterOptionManager IFilterOptionHandlerManager
}

func NewModuleBuilder(engine *gin.Engine, core CoreModule, filterOptionManager ...IFilterOptionHandlerManager) *ModuleBuilder {

	routes := engine.Routes()
	if len(routes) > 0 {
		for _, r := range routes {
			log.DebugF("unknown %s %s %s", r.Handler, r.Method, r.Path)
		}
		panic("unknown")
	}
	b := &ModuleBuilder{
		engine:          engine,
		modules:         make(map[string]Module),
		middlewares:     nil,
		providerBuilder: NewProviderBuilder(),
		core:            core,
		coreMiddleware:  nil,
	}
	if len(filterOptionManager) > 0 {
		b.filterOptionManager = filterOptionManager[0]
	} else {
		b.filterOptionManager = nil
	}
	if middleware, has := core.Middleware(); has {
		midCore := middleware.MiddlewaresInfo()
		mds := make([]MiddlewareHandler, 0, len(systemMiddlewares)+len(midCore))
		if len(systemMiddlewares) > 0 {
			mds = append(mds, systemMiddlewares...)
		}
		mds = append(mds, midCore...)
		b.coreMiddleware = mds

	}
	if providerSupport, has := core.Support(); has {
		b.providerBuilder.Append(core.Name(), providerSupport)
	}
	return b
}

func (b *ModuleBuilder) Append(module ...Module) *ModuleBuilder {
	if len(module) > 0 {

		b.lock.Lock()
		defer b.lock.Unlock()
		for _, m := range module {
			moduleName := m.Name()
			b.modules[moduleName] = m

			if middleware, has := m.Middleware(); has {
				for _, mid := range middleware.MiddlewaresInfo() {
					b.middlewaresAppend = append(b.middlewaresAppend, mid)
				}
			}
			if providerSupport, has := m.Support(); has {
				b.providerBuilder.Append(moduleName, providerSupport)
			}

		}

	}
	return b
}
func (b *ModuleBuilder) createMiddleware() []MiddlewareHandler {
	cms := make(map[string]struct{})
	for _, m := range b.coreMiddleware {
		if m.Alternative {
			cms[m.Name] = struct{}{}
		}
	}
	for _, m := range b.middlewaresAppend {
		delete(cms, m.Name)
	}
	rs := make([]MiddlewareHandler, 0, len(cms)+len(b.middlewaresAppend))
	for _, m := range b.coreMiddleware {
		if _, has := cms[m.Name]; has {
			rs = append(rs, m)
		}
	}
	rs = append(rs, b.middlewaresAppend...)
	return rs
}
func (b *ModuleBuilder) initCore() error {
	if b.core == nil {
		return nil
	}
	routers, ok := b.core.Routers()
	if !ok {
		return errors.New("core must routers")
	}
	routersInfo := routers.RoutersInfo()
	b.alternativeRouters = make(map[string]RouterInfo, len(routersInfo))
	for _, r := range routersInfo {
		if r.Alternative {
			b.alternativeRouters[fmt.Sprintf("%s:%s", r.Method, r.Path)] = r
		}
	}
	b.middlewares = b.createMiddleware()
	return nil
}

// Build 构建模块路由、providers等
func (b *ModuleBuilder) Build() (httpHandler http.Handler, iproviders IProviders, errResult error) {
	defer func() {
		if e := recover(); e != nil {
			errResult = fmt.Errorf("%v", e)
		}
	}()
	b.lock.Lock()
	defer b.lock.Unlock()

	err := b.initCore()
	if err != nil {
		return nil, nil, err
	}

	routerInfo, _ := b.core.Routers()
	log.Debug("handler: core")
	b.handleModule(routerInfo.RoutersInfo())
	for _, m := range b.modules {
		if rs, has := m.Routers(); has {
			log.Debug("handler: ", m.Name())

			errHandler := b.handleModule(rs.RoutersInfo(), ModuleNameHandler(m.Name()))
			if errHandler != nil {
				return nil, nil, errHandler
			}
		}
	}
	log.Debug("handler: alternativeRouters")
	for _, r := range b.alternativeRouters {
		errCore := b.handleRouter(r)
		if errCore != nil {
			return nil, nil, errCore
		}
	}
	b.resetFiltersHandler()
	return b.engine.Handler(), b.providerBuilder.Build(), nil
}
func (b *ModuleBuilder) handleModule(routers RoutersInfo, moduleNameHandler ...gin.HandlerFunc) (errResult error) {

	for i, routerInfo := range routers {
		if routerInfo.Alternative {
			log.DebugF("skip:%s %s %d\n", routerInfo.Method, routerInfo.Path, i)
			continue
		}

		key := fmt.Sprintf("%s:%s", routerInfo.Method, routerInfo.Path)
		if d, has := b.alternativeRouters[key]; has {
			delete(b.alternativeRouters, key)
			routerInfo.Labels = d.Labels
		}
		err := b.handleRouter(routerInfo, moduleNameHandler...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *ModuleBuilder) handleRouter(routerInfo RouterInfo, moduleNameHandler ...gin.HandlerFunc) (errResult error) {
	defer func() {
		if e := recover(); e != nil {
			errResult = fmt.Errorf("%v for handler %s:%s %s ", e, routerInfo.Handler, routerInfo.Method, routerInfo.Path)
		}
	}()
	hs := make([]HandlerFunc, 0, len(routerInfo.HandlerFunc)+1+len(moduleNameHandler))
	hs = append(hs, moduleNameHandler...)
	for _, g := range b.middlewares {
		if g.checkRouter(routerInfo) {
			hs = append(hs, g.Handler)
		}
	}
	hs = append(hs, routerInfo.HandlerFunc...)
	log.DebugF("register:%s %s\n", routerInfo.Method, routerInfo.Path)

	b.engine.Handle(routerInfo.Method, routerInfo.Path, hs...)
	return nil
}

func (b *ModuleBuilder) resetFiltersHandler() {
	if b.filterOptionManager == nil {
		return
	}
	handlers := make(map[string]IFilterOptionHandler)

	if hs, ok := b.core.(IFilterOptionHandlerSupport); ok {
		for _, h := range hs.FilterOptionHandler() {
			handlers[h.Name()] = h
		}
	}
	for _, m := range b.modules {
		if hs, ok := m.(IFilterOptionHandlerSupport); ok {
			for _, h := range hs.FilterOptionHandler() {
				handlers[h.Name()] = h
			}
		}
	}

	b.filterOptionManager.ResetFilterOptionHandlers(handlers)

}
