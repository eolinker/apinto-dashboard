package service

import (
	"context"
	"github.com/eolinker/apinto-dashboard/cache"
	"github.com/eolinker/apinto-dashboard/modules/core"
	module_plugin "github.com/eolinker/apinto-dashboard/modules/module-plugin"
	apinto_module "github.com/eolinker/apinto-module"
	"github.com/eolinker/eosc/common/bean"
	"github.com/eolinker/eosc/log"
	"github.com/go-basic/uuid"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

var (
	_ core.ICore = (*coreService)(nil)
)

const (
	moduleConfigVersionKey = "apinto.module:version"
)

type tModulesData struct {
	data map[string]struct{}
}

func newTModulesData() *tModulesData {
	return &tModulesData{data: map[string]struct{}{}}
}

type coreService struct {
	handlerPointer atomic.Pointer[http.Handler]
	localVersion   string
	lock           sync.Mutex

	modulePluginService module_plugin.IModulePlugin
	engineCreate        core.EngineCreate
	providerService     IProviderService
	modulesData         *tModulesData
	cacheCommon         cache.ICommonCache
	once                sync.Once
}

func (c *coreService) HasModule(module string, path string) bool {
	if c.modulesData == nil {
		return false
	}
	_, has := c.modulesData.data[module]
	return has
}

func (c *coreService) CheckNewModule(name, driver string, config, define interface{}) error {
	return nil
}

func (c *coreService) ResetVersion(version string) {
	if version == "" {
		version = uuid.New()
	}
	c.lock.Lock()
	defer c.lock.Unlock()

	err := c.cacheCommon.Set(context.Background(), moduleConfigVersionKey, []byte(version), 0)
	go c.reloadModule(version)
	if err != nil {
		log.Errorf("set module config version:%s", err.Error())

		return
	}

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
	c.once.Do(func() {
		c.cacheCommon.SetNX(context.Background(), moduleConfigVersionKey, uuid.New(), 0)
		go c.doLoop()
	})
	version, err := c.cacheCommon.Get(context.Background(), moduleConfigVersionKey)
	if err != nil {
		log.Errorf("get module config version:%s", err.Error())

		return err
	}

	return c.reloadModule(string(version))
}
func (c *coreService) doLoop() {
	tick := time.Tick(time.Second * 10)
	for range tick {
		version, err := c.cacheCommon.Get(context.Background(), moduleConfigVersionKey)
		if err != nil {
			log.Errorf("get module config version:%s", err.Error())
			continue
		}
		c.reloadModule(string(version))
	}

}
func (c *coreService) reloadModule(version string) error {

	c.lock.Lock()
	defer c.lock.Unlock()

	if c.localVersion != version {
		c.localVersion = version
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
	modulesData := newTModulesData()
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
		modulesData.data[module.Name] = struct{}{}
		builder.Append(m)
	}

	handler, provider, err := builder.Build()
	if err != nil {
		return err
	}
	c.handlerPointer.Store(&handler)
	c.providerService.set(provider)
	c.modulesData = modulesData
	return nil
}
func NewService(providerService IProviderService) core.ICore {

	c := &coreService{
		providerService: providerService,
	}
	bean.Autowired(&c.modulePluginService)
	bean.Autowired(&c.engineCreate)
	bean.Autowired(&c.cacheCommon)
	return c
}
