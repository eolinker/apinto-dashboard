package apinto_module

import (
	"fmt"
	"github.com/eolinker/eosc/log"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

var (
	systemModules []Module
)

func AddSystemModule(m Module) {
	systemModules = append(systemModules, m)
}

type ModuleBuilder struct {
	lock    sync.Mutex
	modules map[string]Module

	engine *gin.Engine

	filterOptionManager IFilterOptionHandlerManager
}
type builderHandler struct {
	ReplaceableRouters map[string]RouterInfo
	middlewares        []MiddlewareHandler
	engine             *gin.Engine
	providerBuilder    *ProviderBuilder
}

func (b *ModuleBuilder) newBuilderHandler() *builderHandler {
	middlewares, providerBuilder, replaceableRouters := initSystem(b.modules)
	hs := &builderHandler{
		engine:             b.engine,
		middlewares:        middlewares,
		ReplaceableRouters: replaceableRouters,
		providerBuilder:    providerBuilder,
	}

	return hs
}

func initSystem(modules map[string]Module) ([]MiddlewareHandler, *ProviderBuilder, map[string]RouterInfo) {
	provider := NewProviderBuilder()
	systemMiddlewares := make([]MiddlewareHandler, 0)
	cms := make(map[string]struct{})
	ReplaceableRouters := make(map[string]RouterInfo)

	for _, m := range systemModules {
		if rs, ok := m.Routers(); ok {
			for _, r := range rs.RoutersInfo() {
				if r.Replaceable {
					ReplaceableRouters[fmt.Sprintf("%s:%s", r.Method, r.Path)] = r
				}
			}
		}
		if md, ok := m.Middleware(); ok {

			for _, mdl := range md.MiddlewaresInfo() {
				if mdl.Replaceable {
					cms[mdl.Name] = struct{}{}
				}
				systemMiddlewares = append(systemMiddlewares, mdl)
			}
		}
		if p, ok := m.Support(); ok {
			provider.Append(m.Name(), p)
		}
	}

	mdappend := make([]MiddlewareHandler, 0)
	for _, m := range modules {
		if md, ok := m.Middleware(); ok {
			for _, mdh := range md.MiddlewaresInfo() {
				delete(cms, mdh.Name)
				mdappend = append(mdappend, mdh)
			}
		}
		if p, ok := m.Support(); ok {
			provider.Append(m.Name(), p)
		}
	}
	mdws := make([]MiddlewareHandler, 0, len(cms)+len(mdappend))
	for _, sm := range systemMiddlewares {
		if !sm.Replaceable {
			mdws = append(mdws, sm)
			continue
		}
		if _, has := cms[sm.Name]; has {
			mdws = append(mdws, sm)
		}
	}
	mdws = append(mdws, mdappend...)
	return mdws, provider, ReplaceableRouters
}
func NewModuleBuilder(engine *gin.Engine, filterOptionManager ...IFilterOptionHandlerManager) *ModuleBuilder {

	routes := engine.Routes()
	if len(routes) > 0 {
		for _, r := range routes {
			log.DebugF("unknown %s %s %s", r.Handler, r.Method, r.Path)
		}
		panic("unknown")
	}
	b := &ModuleBuilder{
		engine:  engine,
		modules: make(map[string]Module),
	}
	if len(filterOptionManager) > 0 {
		b.filterOptionManager = filterOptionManager[0]
	} else {
		b.filterOptionManager = nil
	}
	//if middleware, has := core.Middleware(); has {
	//	midCore := middleware.MiddlewaresInfo()
	//	mds := make([]MiddlewareHandler, 0, len(systemMiddlewares)+len(midCore))
	//	if len(systemMiddlewares) > 0 {
	//		mds = append(mds, systemMiddlewares...)
	//	}
	//	mds = append(mds, midCore...)
	//	b.coreMiddleware = mds
	//
	//}
	//if providerSupport, has := core.Support(); has {
	//	b.providerBuilder.Append(core.Name(), providerSupport)
	//}
	return b
}

func (b *ModuleBuilder) Append(module ...Module) *ModuleBuilder {
	if len(module) > 0 {

		b.lock.Lock()
		defer b.lock.Unlock()
		for _, m := range module {
			moduleName := m.Name()
			b.modules[moduleName] = m

			//if middleware, has := m.Middleware(); has {
			//	for _, mid := range middleware.MiddlewaresInfo() {
			//		b.middlewaresAppend = append(b.middlewaresAppend, mid)
			//	}
			//}
			//if providerSupport, has := m.Support(); has {
			//	b.providerBuilder.Append(moduleName, providerSupport)
			//}

		}

	}
	return b
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

	bh := b.newBuilderHandler()
	for _, m := range systemModules {

		if rs, ok := m.Routers(); ok {
			log.Debug("handler: ", m.Name())
			err := bh.handleModule(rs.RoutersInfo(), ModuleNameHandler(m.Name()))
			if err != nil {
				return nil, nil, err
			}
		}
	}

	for _, m := range b.modules {
		if rs, has := m.Routers(); has {
			log.Debug("handler: ", m.Name())

			errHandler := bh.handleModule(rs.RoutersInfo(), ModuleNameHandler(m.Name()))
			if errHandler != nil {
				return nil, nil, errHandler
			}
		}
	}
	log.Debug("handler: ReplaceableRouters")
	for _, r := range bh.ReplaceableRouters {
		errCore := bh.handleRouter(r)
		if errCore != nil {
			return nil, nil, errCore
		}
	}
	b.resetFiltersHandler()
	return bh.engine.Handler(), bh.providerBuilder.Build(), nil
}
func (b *builderHandler) handleModule(routers RoutersInfo, moduleNameHandler ...gin.HandlerFunc) (errResult error) {

	for i, routerInfo := range routers {
		if routerInfo.Replaceable {
			log.DebugF("skip:%s %s %d\n", routerInfo.Method, routerInfo.Path, i)
			continue
		}

		key := fmt.Sprintf("%s:%s", routerInfo.Method, routerInfo.Path)
		if d, has := b.ReplaceableRouters[key]; has {
			delete(b.ReplaceableRouters, key)
			routerInfo.Labels = d.Labels
		}
		err := b.handleRouter(routerInfo, moduleNameHandler...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *builderHandler) handleRouter(routerInfo RouterInfo, moduleNameHandler ...gin.HandlerFunc) (errResult error) {
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

	for _, m := range systemModules {
		if hs, ok := m.(IFilterOptionHandlerSupport); ok {
			for _, h := range hs.FilterOptionHandler() {
				handlers[h.Name()] = h
			}
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
