package dynamic_controller

import (
	"context"
	"fmt"
	v2 "github.com/eolinker/apinto-dashboard/client/v2"
	"github.com/eolinker/apinto-dashboard/modules/dynamic"
	"github.com/eolinker/eosc/common/bean"
	"github.com/eolinker/eosc/log"
	"net/http"

	"github.com/eolinker/apinto-module"
)

type DynamicModuleDriver struct {
}

func NewDynamicModuleDriver() apinto_module.Driver {
	return &DynamicModuleDriver{}
}

func (c *DynamicModuleDriver) CreatePlugin(define interface{}) (apinto_module.Plugin, error) {
	return &DynamicModulePlugin{
		define: define,
	}, nil
}

func (c *DynamicModulePlugin) GetPluginFrontend(moduleName string) string {
	return fmt.Sprintf("template/%s", moduleName)
}

func (c *DynamicModulePlugin) IsPluginVisible() bool {
	return true
}

func (c *DynamicModulePlugin) IsShowServer() bool {
	return false
}

func (c *DynamicModulePlugin) IsCanUninstall() bool {
	return true
}

func (c *DynamicModulePlugin) IsCanDisable() bool {
	return true
}

type DynamicModulePlugin struct {
	define interface{}
}

func (d *DynamicModulePlugin) CreateModule(name string, config interface{}) (apinto_module.Module, error) {
	return NewDynamicModule(name, d.define), nil
}

func (d *DynamicModulePlugin) CheckConfig(name string, config interface{}) error {
	return nil
}

type DynamicModule struct {
	name           string
	define         interface{}
	routers        apinto_module.RoutersInfo
	profession     string
	skill          string
	dynamicService dynamic.IDynamicService
}

func (c *DynamicModule) Provider() map[string]apinto_module.Provider {
	return map[string]apinto_module.Provider{
		c.skill: newSkillProvider(c.profession, c.skill),
	}
}
func (c *DynamicModule) Status(key string, namespaceId int, cluster string) apinto_module.CargoStatus {
	status, err := c.dynamicService.ClusterStatusByClusterName(context.Background(), namespaceId, c.profession, key, cluster)
	if err != nil {
		log.Error(err)
		return apinto_module.None
	}
	if status.Status == v2.StatusOnline || status.Status == v2.StatusPre {
		return apinto_module.Online
	}
	return apinto_module.Offline

}

func (c *DynamicModule) Name() string {
	return c.name
}

func (c *DynamicModule) Support() (apinto_module.ProviderSupport, bool) {
	return c, true
}

func (c *DynamicModule) Routers() (apinto_module.Routers, bool) {
	return c, true
}

func (c *DynamicModule) Middleware() (apinto_module.Middleware, bool) {
	return nil, false
}

func NewDynamicModule(name string, define interface{}) *DynamicModule {
	dm := &DynamicModule{name: name, define: define}
	bean.Autowired(&dm.dynamicService)
	dm.initRouter()
	return dm
}

func (c *DynamicModule) RoutersInfo() apinto_module.RoutersInfo {
	return c.routers
}
func (c *DynamicModule) initRouter() {
	dc := newDynamicController(c.name, c.define)
	c.profession = dc.Profession
	c.skill = dc.Skill
	c.routers = []apinto_module.RouterInfo{
		{
			Method:      http.MethodGet,
			Path:        fmt.Sprintf("/api/dynamic/%s/list", c.name),
			Handler:     "dynamic.list",
			HandlerFunc: []apinto_module.HandlerFunc{dc.list},
		},
		{
			Method:      http.MethodGet,
			Path:        fmt.Sprintf("/api/dynamic/%s/info/:uuid", c.name),
			Handler:     "dynamic.info",
			HandlerFunc: []apinto_module.HandlerFunc{dc.info},
		},
		{
			Method:      http.MethodGet,
			Path:        fmt.Sprintf("/api/dynamic/%s/render", c.name),
			Handler:     "dynamic.render",
			HandlerFunc: []apinto_module.HandlerFunc{dc.render},
		},
		{
			Method:      http.MethodGet,
			Path:        fmt.Sprintf("/api/dynamic/%s/cluster/:uuid", c.name),
			Handler:     "dynamic.cluster_status",
			HandlerFunc: []apinto_module.HandlerFunc{dc.clusterStatus},
		},
		{
			Method:      http.MethodGet,
			Path:        fmt.Sprintf("/api/dynamic/%s/status", c.name),
			Handler:     "dynamic.cluster_statuses",
			HandlerFunc: []apinto_module.HandlerFunc{dc.clusterStatusList},
		},
		{
			Method:      http.MethodDelete,
			Path:        fmt.Sprintf("/api/dynamic/%s/batch", c.name),
			Handler:     "dynamic.delete",
			HandlerFunc: []apinto_module.HandlerFunc{dc.batchDelete},
		}, {
			Method:      http.MethodPost,
			Path:        fmt.Sprintf("/api/dynamic/%s", c.name),
			Handler:     "dynamic.save",
			HandlerFunc: []apinto_module.HandlerFunc{dc.create},
		}, {
			Method:      http.MethodPut,
			Path:        fmt.Sprintf("/api/dynamic/%s/config/:uuid", c.name),
			Handler:     "dynamic.save",
			HandlerFunc: []apinto_module.HandlerFunc{dc.save},
		},
		{
			Method:      http.MethodPut,
			Path:        fmt.Sprintf("/api/dynamic/%s/online/:uuid", c.name),
			Handler:     "dynamic.online",
			HandlerFunc: []apinto_module.HandlerFunc{dc.online},
		},
		{
			Method:      http.MethodPut,
			Path:        fmt.Sprintf("/api/dynamic/%s/offline/:uuid", c.name),
			Handler:     "dynamic.offline",
			HandlerFunc: []apinto_module.HandlerFunc{dc.offline},
		},
	}
}
