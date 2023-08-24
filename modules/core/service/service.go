package service

import (
	"context"
	"fmt"
	"github.com/eolinker/apinto-dashboard/cache"
	apinto_module "github.com/eolinker/apinto-dashboard/module"
	"github.com/eolinker/apinto-dashboard/modules/core"
	"github.com/eolinker/apinto-dashboard/modules/core/controller"
	module_plugin "github.com/eolinker/apinto-dashboard/modules/module-plugin"
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
	data map[string]apinto_module.Module
}

func newTModulesData() *tModulesData {
	return &tModulesData{data: map[string]apinto_module.Module{}}
}

type coreService struct {
	handlerPointer atomic.Pointer[http.Handler]
	localVersion   string
	lock           sync.RWMutex

	modulePluginService module_plugin.IModulePlugin
	engineCreate        core.EngineCreate
	providerService     IProviderService
	modulesData         *tModulesData
	cacheCommon         cache.ICommonCache
	once                sync.Once

	coreModule apinto_module.CoreModule

	filterOptionHandlerManager apinto_module.IFilterOptionHandlerManager
}

func (c *coreService) SetCoreModule(module apinto_module.CoreModule) {
	c.coreModule = module
}

func (c *coreService) HasModule(module string, path string) bool {
	if c.modulesData == nil {
		return false
	}
	c.lock.RLock()
	defer c.lock.RUnlock()
	_, has := c.modulesData.data[module]
	return has
}

func (c *coreService) CheckNewModule(UUID, name, driverName string, define, config interface{}) error {
	_, err := createModule(driverName, name, define, config)
	if err != nil {
		return err
	}
	c.lock.Lock()
	defer c.lock.Unlock()

	ctx := context.Background()
	modules, err := c.modulePluginService.GetEnabledPlugins(ctx)
	if err != nil {
		return err
	}
	//builder := apinto_module.NewModuleBuilder(c.engineCreate.CreateEngine())
	for _, module := range modules {
		if module.Name == name {
			if module.UUID == UUID {
				continue
			} else {
				return fmt.Errorf("%w on %s", apinto_module.ErrorModuleNameConflict, module.UUID)
			}
		}
		//m, err := createModule(module.Driver, module.Name, module.Define, module.Config)
		//if err != nil {
		//	continue
		//}
		//builder.Append(m)
	}
	return nil
	//builder.Append(newModule)
	//
	//_, _, err = builder.Build()
	//return err

}
func createModule(driverName, name string, define, config interface{}) (apinto_module.Module, error) {
	driver, has := apinto_module.GetDriver(driverName)
	if !has {

		err := fmt.Errorf("%w %s", apinto_module.ErrorDriverNotExist, driverName)
		log.Error(err)
		return nil, err
	}
	plugin, err := driver.CreatePlugin(define)
	if err != nil {
		err2 := fmt.Errorf("create plugin %s error:%w", name, err)
		log.Error(err2)
		return nil, err2
	}
	err = plugin.CheckConfig(name, config)
	if err != nil {

		err2 := fmt.Errorf("plugin module %s config error:%w", name, err)
		log.Error(err2)
		return nil, err2
	}

	m, err := plugin.CreateModule(name, config)
	if err != nil {

		err2 := fmt.Errorf("create module %s  error:%w", name, err)
		log.Error(err2)
		return nil, err2
	}
	return m, nil
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
	tick := time.NewTicker(time.Second * 10)
	defer tick.Stop()
	for range tick.C {
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
	builder := apinto_module.NewModuleBuilder(c.engineCreate.CreateEngine(), c.filterOptionHandlerManager)
	for _, module := range modules {
		m, err := createModule(module.Driver, module.Name, module.Define, module.Config)
		if err != nil {
			log.Errorf("create module %s  error:%s", module.Name, err.Error())
			continue
		}

		modulesData.data[module.Name] = m
		builder.Append(m)
	}

	handler, provider, err := builder.Build()
	if err != nil {
		return err
	}
	oldModules := c.modulesData
	defer func() {
		if oldModules == nil {
			return
		}
		for _, m := range oldModules.data {
			if k, ok := m.(apinto_module.ModuleNeedKill); ok {
				k.Kill()
			}
		}
	}()
	c.handlerPointer.Store(&handler)
	c.providerService.set(provider)
	c.modulesData = modulesData
	return nil
}
func NewService(providerService IProviderService) core.ICore {

	c := &coreService{
		providerService: providerService,
		coreModule:      controller.NewModule(),
	}
	apinto_module.AddSystemModule(c.coreModule)
	bean.Autowired(&c.filterOptionHandlerManager)
	bean.Autowired(&c.modulePluginService)
	bean.Autowired(&c.engineCreate)
	bean.Autowired(&c.cacheCommon)
	return c
}
