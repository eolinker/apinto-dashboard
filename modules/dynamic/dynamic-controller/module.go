package dynamic_controller

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/eolinker/apinto-dashboard/common"

	v2 "github.com/eolinker/apinto-dashboard/client/v2"
	"github.com/eolinker/apinto-dashboard/modules/dynamic"
	"github.com/eolinker/eosc/common/bean"
	"github.com/eolinker/eosc/log"

	"github.com/eolinker/apinto-dashboard/module"
)

type DynamicModuleDriver struct {
	showServer bool
}

func NewDynamicModuleDriver(showServer bool) *DynamicModuleDriver {
	return &DynamicModuleDriver{showServer: showServer}
}

func (c *DynamicModuleDriver) CreatePlugin(define interface{}) (apinto_module.Plugin, error) {
	v, err := apinto_module.DecodeFor[DynamicDefine](define)
	if err != nil {
		return nil, err
	}
	return &DynamicModulePlugin{
		DynamicModuleDriver: c,
		define:              v,
	}, nil
}

func (c *DynamicModulePlugin) GetPluginFrontend(moduleName string) string {
	return fmt.Sprintf("template/%s", moduleName)
}

func (c *DynamicModulePlugin) IsShowServer() bool {
	return c.showServer
}

type DynamicModulePlugin struct {
	*DynamicModuleDriver
	define *DynamicDefine
}

func (d *DynamicModulePlugin) CreateModule(name string, config interface{}) (apinto_module.Module, error) {

	return NewDynamicModule(name, d.define), nil
}

func (d *DynamicModulePlugin) CheckConfig(name string, config interface{}) error {
	return nil
}

type DynamicModule struct {
	name                 string
	define               *DynamicDefine
	routers              apinto_module.RoutersInfo
	profession           string
	skill                string
	dynamicService       dynamic.IDynamicService
	filterOptionHandlers []apinto_module.IFilterOptionHandler
}

func (c *DynamicModule) FilterOptionHandler() []apinto_module.IFilterOptionHandler {
	return c.filterOptionHandlers
}

func (c *DynamicModule) Provider() map[string]apinto_module.Provider {
	return map[string]apinto_module.Provider{
		c.skill: newSkillProvider(c.profession, c.skill),
	}
}

func (c *DynamicModule) Status(key string, namespaceId int, cluster string) (apinto_module.CargoStatus, string) {
	id := fmt.Sprintf("%s@%s", strings.ToLower(key), c.profession)
	if cluster == "" {
		return apinto_module.None, id
	}
	status, err := c.dynamicService.ClusterStatusByClusterName(context.Background(), namespaceId, c.profession, key, cluster)
	if err != nil {
		log.Error(err)
		return apinto_module.None, id
	}

	if status.Status == v2.StatusOnline || status.Status == v2.StatusPre {
		return apinto_module.Online, id
	}
	return apinto_module.Offline, id

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

func NewDynamicModule(name string, define *DynamicDefine) *DynamicModule {
	dm := &DynamicModule{name: name, define: define}

	bean.Autowired(&dm.dynamicService)
	dm.initRouter()
	dm.initFilter()
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
func (c *DynamicModule) initFilter() {
	if c.define.FilterOptions == nil {
		return
	}
	cf := c.define.FilterOptions
	c.filterOptionHandlers = []apinto_module.IFilterOptionHandler{
		NewFilterOption(cf.Name, apinto_module.FilterOptionConfig{

			Title: cf.Title,
			Titles: common.SliceToSlice(cf.Titles, func(s OptionTitle) apinto_module.OptionTitle {
				return apinto_module.OptionTitle{
					Title: s.Title,
					Field: s.Field,
				}
			}),
			Key: "id",
		}, common.SliceToSlice(c.define.Fields, func(s *Basic) string {
			return s.Name
		}), c.profession, common.SliceToSlice(c.define.Drivers, func(s *Basic) string {
			return s.Name
		})),
	}
}
