package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eolinker/apinto-dashboard/cache"
	apinto_module "github.com/eolinker/apinto-dashboard/module"
	"github.com/eolinker/apinto-dashboard/module/builder"
	"github.com/eolinker/apinto-dashboard/modules/core"
	"github.com/eolinker/apinto-dashboard/modules/mpm3"
	"github.com/eolinker/apinto-dashboard/pm3"
	"github.com/eolinker/eosc/common/bean"
	"github.com/eolinker/eosc/log"
	"github.com/go-basic/uuid"
	"github.com/redis/go-redis/v9"
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

	pluginService   mpm3.IPluginService
	engineCreate    core.EngineCreate
	providerService IProviderService
	modulesData     *tModulesData
	cacheCommon     cache.ICommonCache
	once            sync.Once

	filterOptionHandlerManager apinto_module.IFilterOptionHandlerManager
}

//	func (c *coreService) CheckPluginConfig(UUID, name, driverName string, define *pm3.PluginDefine, config []byte) error {
//		newModule, err := createModule(driverName, define, config)
//		if err != nil {
//			return err
//		}
//		c.lock.Lock()
//		defer c.lock.Unlock()
//
//		ctx := context.Background()
//		modules, err := c.pluginService.GetPluginConfig(ctx)
//		if err != nil {
//			return err
//		}
//		builder := builder.NewModuleBuilder(c.engineCreate.CreateEngine())
//		for _, module := range modules {
//			if module.UUID == UUID {
//				continue
//			}
//			m, err := createModule(module.Driver, module.Define, module.Config)
//			if err != nil {
//				continue
//			}
//			builder.Append(m)
//		}
//
//		builder.Append(newModule)
//
//		_, _, err = builder.Build()
//		//return err
//
// }
func createModule(driverName string, define *pm3.PluginDefine, config pm3.PluginConfig) (pm3.Module, error) {
	driver, has := apinto_module.GetDriver(driverName)
	if !has {

		err := fmt.Errorf("%w %s", apinto_module.ErrorDriverNotExist, driverName)
		log.Error(err)
		return nil, err
	}

	plugin, err := driver.Create(define, config)
	if err != nil {
		err2 := fmt.Errorf("create plugin %s error:%w", define.Id, err)
		log.Error(err2)
		return nil, err2
	}

	return plugin, nil
}
func (c *coreService) ResetVersion(version string) {

	c.lock.Lock()
	defer c.lock.Unlock()
	if version == "" {
		if c.localVersion != "" {
			return
		}
		version = uuid.New()
	}
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
		go c.doLoop()
	})
	version, err := c.cacheCommon.Get(context.Background(), moduleConfigVersionKey)
	if err != nil {
		log.Errorf("get module config version:%s", err.Error())
	}

	c.ResetVersion(string(version))
	return nil
}
func (c *coreService) doLoop() {
	tick := time.NewTicker(time.Second * 10)
	defer tick.Stop()
	ctx := context.Background()
	for range tick.C {
		timeout, _ := context.WithTimeout(ctx, time.Second)
		version, err := c.cacheCommon.Get(timeout, moduleConfigVersionKey)
		if err != nil {
			if errors.Is(err, redis.Nil) {
				c.ResetVersion("")

			} else {
				log.Errorf("get module config version:%s", err.Error())
			}
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
	plugins, err := c.pluginService.GetEnabled(ctx)
	if err != nil {
		return err
	}

	modulesData := newTModulesData()
	builder := builder.NewModuleBuilder(c.engineCreate.CreateEngine(), c.filterOptionHandlerManager)
	for _, module := range plugins {
		var cf pm3.PluginConfig
		if len(module.Config) > 0 {
			json.Unmarshal(module.Config, &cf)
		}
		m, err := createModule(module.Driver, module.Define, cf)
		if err != nil {
			log.Errorf("create module %s  error:%s", module.Define.Id, err.Error())
			continue
		}

		modulesData.data[module.UUID] = m
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
	}

	bean.Autowired(&c.filterOptionHandlerManager)
	bean.Autowired(&c.pluginService)
	bean.Autowired(&c.engineCreate)
	bean.Autowired(&c.cacheCommon)

	return c
}
