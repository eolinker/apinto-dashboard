package service

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/core"
	"github.com/eolinker/apinto-dashboard/modules/middleware"
	module_plugin "github.com/eolinker/apinto-dashboard/modules/module-plugin"
	apinto_module "github.com/eolinker/apinto-module"
	"github.com/eolinker/eosc/common/bean"
	"github.com/eolinker/eosc/log"
	"net/http"
	"sync"
	"sync/atomic"
)

var (
	_ core.ICore = (*coreService)(nil)
)

type coreService struct {
	handlerPointer      atomic.Pointer[http.Handler]
	localVersionPointer atomic.Pointer[string]
	lock                sync.Mutex

	middlewareService   middleware.IMiddlewareService
	modulePluginService module_plugin.IModulePluginService
	engineCreate        core.EngineCreate
	providerService     core.IProviders
}

func (c *coreService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := c.handlerPointer.Load()
	if handler == nil || (*handler) == nil {
		http.NotFound(w, r)
		return
	}
	(*handler).ServeHTTP(w, r)
}

func (c *coreService) ReloadModule(version string) error {

	localVersion := c.localVersionPointer.Swap(&version)

	if localVersion != nil && (*localVersion) == version {
		return nil
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	localVersion = c.localVersionPointer.Load()
	if localVersion != nil && (*localVersion) == version {
		// todo load module
		// todo load middleware
		err := c.rebuild()
		if err != nil {
			log.Error("error to rebuild core:", err)
			return err
		}
	}
	return nil
}
func (c *coreService) rebuild() error {
	ctx := context.Background()
	modules, err := c.modulePluginService.GetEnabledPlugins(ctx)
	if err != nil {
		return err
	}
	middlewaresConfig, err := c.middlewareService.Groups(ctx)
	if err != nil {
		return err
	}
	builder := apinto_module.NewModuleBuilder(c.engineCreate.CreateEngine())
	for _, module := range modules {
		driver, has := apinto_module.GetDriver(module.Driver)
		if !has {
			log.Errorf("not find driver:%s", module.Driver)
			continue
		}
		plugin, err := driver.CreatePlugin(module.Define)
		if err != nil {
			log.Errorf("create plugin %s error:%s", module.Name, err.Error())
			continue
		}
		err = plugin.CheckConfig(module.Name, module.Config.APIGroup, module.Config)
		if err != nil {
			log.Errorf("plugin module %s config error:%s", module.Name, err.Error())
			continue
		}

		m, err := plugin.CreateModule(module.Name, module.Config.APIGroup, module.Config)
		if err != nil {
			log.Errorf("create module %s  error:%s", module.Name, err.Error())
			continue
		}
		builder.Append(m)
	}
	middlewares := make([]apinto_module.MiddlewareConfig, 0, len(middlewaresConfig))
	for _, mid := range middlewaresConfig {
		middlewares = append(middlewares, apinto_module.MiddlewareConfig{
			Prefix:   mid.Prefix,
			Handlers: mid.Middlewares,
		})
	}
	handler, provider, err := builder.Build(middlewares)
	if err != nil {
		return err
	}
	c.handlerPointer.Store(&handler)
	c.providerService.Set(provider)
	return nil
}
func NewService(providerService core.IProviders) core.ICore {

	c := &coreService{
		providerService: providerService,
	}
	bean.Autowired(&c.modulePluginService)
	bean.Autowired(&c.middlewareService)
	bean.Autowired(&c.engineCreate)
	bean.AddInitializingBeanFunc(func() {
		c.rebuild()
	})
	return c
}
