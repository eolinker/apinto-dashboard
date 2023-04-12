package service

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/core"
	module_plugin "github.com/eolinker/apinto-dashboard/modules/module-plugin"
	apinto_module "github.com/eolinker/apinto-module"
	"github.com/eolinker/eosc/common/bean"
	"github.com/eolinker/eosc/log"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var (
	_ core.ICore = (*coreService)(nil)
)

type coreService struct {
	handlerPointer atomic.Pointer[http.Handler]
	localVersion   string
	lock           sync.Mutex

	modulePluginService module_plugin.IModulePlugin
	engineCreate        core.EngineCreate
	providerService     IProviderService
}

func (c *coreService) HasModule(module string, path string) bool {
	return true
}

func (c *coreService) CheckNewModule(pluginID, name string, config interface{}) error {
	return nil
}

func (c *coreService) ResetVersion(version string) {

}

func (c *coreService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := c.handlerPointer.Load()
	if handler == nil || (*handler) == nil {
		http.NotFound(w, r)
		return
	}
	(*handler).ServeHTTP(w, r)
}

func (c *coreService) ReloadModule() error {

	lastVersion := strconv.FormatInt(time.Now().UnixNano(), 16) // todo load lastVersion from redis or db
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.localVersion != lastVersion {
		c.localVersion = lastVersion
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
		err = plugin.CheckConfig(module.Name, module.Config)
		if err != nil {
			log.Errorf("plugin module %s config error:%s", module.Name, err.Error())
			continue
		}

		m, err := plugin.CreateModule(module.Name, module.Config)
		if err != nil {
			log.Errorf("create module %s  error:%s", module.Name, err.Error())
			continue
		}
		builder.Append(m)
	}

	handler, provider, err := builder.Build()
	if err != nil {
		return err
	}
	c.handlerPointer.Store(&handler)
	c.providerService.set(provider)
	return nil
}
func NewService(providerService IProviderService) core.ICore {

	c := &coreService{
		providerService: providerService,
	}
	bean.Autowired(&c.modulePluginService)
	bean.Autowired(&c.engineCreate)

	return c
}
