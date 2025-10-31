package dynamic_controller

import (
	"context"
	"fmt"
	"github.com/eolinker/apinto-dashboard/pm3"
	"net/http"
	"strings"

	"github.com/eolinker/apinto-dashboard/common"

	v2 "github.com/eolinker/apinto-dashboard/client/v2"
	"github.com/eolinker/apinto-dashboard/modules/dynamic"
	"github.com/eolinker/eosc/common/bean"
	"github.com/eolinker/eosc/log"

	"github.com/eolinker/apinto-dashboard/module"
)

var (
	_ apinto_module.Driver = (*DynamicModuleDriver)(nil)
)

type DynamicModuleDriver struct {
}

func (d *DynamicModuleDriver) Install(info *pm3.PluginDefine) (ms []pm3.PModule, acs []pm3.PAccess, fs []pm3.PFrontend, err error) {
	return pm3.ReadPluginAssembly(info)
}

func (d *DynamicModuleDriver) Create(info *pm3.PluginDefine, config pm3.PluginConfig) (pm3.Module, error) {

	return d.NewDynamicModule(info)
}

func NewDynamicModuleDriver() *DynamicModuleDriver {
	return &DynamicModuleDriver{}
}

type DynamicModulePlugin struct {
	*DynamicModuleDriver
	define *DynamicDefine
}

type DynamicModule struct {
	*pm3.ModuleTool

	id                   string
	name                 string
	define               *DynamicDefine
	routers              apinto_module.RoutersInfo
	profession           string
	skill                string
	dynamicService       dynamic.IDynamicService
	filterOptionHandlers []apinto_module.IFilterOptionHandler
}

func (c *DynamicModule) Frontend() []pm3.FrontendAsset {
	return nil
}

func (c *DynamicModule) Apis() []pm3.Api {
	return c.routers
}

func (c *DynamicModule) Middleware() []pm3.Middleware {
	return nil
}

func (c *DynamicModule) FilterOptionHandler() []apinto_module.IFilterOptionHandler {
	return c.filterOptionHandlers
}

func (c *DynamicModule) Provider() map[string]pm3.Provider {
	return map[string]pm3.Provider{
		c.skill: newSkillProvider(c.profession, c.skill),
	}
}

func (c *DynamicModule) Status(key string, namespaceId int, cluster string) (pm3.CargoStatus, string) {
	id := fmt.Sprintf("%s@%s", strings.ToLower(key), c.profession)
	if cluster == "" {
		return pm3.None, id
	}
	status, err := c.dynamicService.ClusterStatusByClusterName(context.Background(), namespaceId, c.profession, key, cluster)
	if err != nil {
		log.Error(err)
		return pm3.None, id
	}

	if status.Status == v2.StatusOnline || status.Status == v2.StatusPre {
		return pm3.Online, id
	}
	return pm3.Offline, id

}

func (c *DynamicModule) Name() string {
	return c.name
}

func (c *DynamicModule) Support() (pm3.ProviderSupport, bool) {
	return c, true
}

func (d *DynamicModuleDriver) NewDynamicModule(info *pm3.PluginDefine) (*DynamicModule, error) {
	define, err := apinto_module.DecodeFor[DynamicDefine](info.Define)
	if err != nil {
		return nil, err
	}
	dm := &DynamicModule{ModuleTool: pm3.NewModuleTool(info.Id, info.Name),
		id: info.Id, name: info.Name, define: define}

	bean.Autowired(&dm.dynamicService)
	dm.initRouter(info.Cname, define)
	dm.initFilter()
	return dm, nil
}

func (c *DynamicModule) initRouter(cname string, define *DynamicDefine) {
	dc := newDynamicController(c.id, c.name, cname, define)
	c.profession = dc.Profession
	c.skill = dc.Skill
	c.routers = []apinto_module.RouterInfo{
		{
			Method: http.MethodGet,
			Path:   fmt.Sprintf("/api/dynamic/%s/list", c.name),

			HandlerFunc: dc.list,
		},
		{
			Method: http.MethodGet,
			Path:   fmt.Sprintf("/api/dynamic/%s/info/:uuid", c.name),

			HandlerFunc: dc.info,
		},
		{
			Method: http.MethodGet,
			Path:   fmt.Sprintf("/api/dynamic/%s/render", c.name),

			HandlerFunc: dc.render,
		},
		{
			Method: http.MethodGet,
			Path:   fmt.Sprintf("/api/dynamic/%s/cluster/:uuid", c.name),

			HandlerFunc: dc.clusterStatus,
		},
		{
			Method: http.MethodGet,
			Path:   fmt.Sprintf("/api/dynamic/%s/status", c.name),

			HandlerFunc: dc.clusterStatusList,
		},
		{
			Method: http.MethodDelete,
			Path:   fmt.Sprintf("/api/dynamic/%s/batch", c.name),

			HandlerFunc: dc.batchDelete,
		}, {
			Method: http.MethodPost,
			Path:   fmt.Sprintf("/api/dynamic/%s", c.name),

			HandlerFunc: dc.create,
		}, {
			Method: http.MethodPut,
			Path:   fmt.Sprintf("/api/dynamic/%s/config/:uuid", c.name),

			HandlerFunc: dc.save,
		},
		{
			Method: http.MethodPut,
			Path:   fmt.Sprintf("/api/dynamic/%s/online/:uuid", c.name),

			HandlerFunc: dc.online,
		},
		{
			Method: http.MethodPut,
			Path:   fmt.Sprintf("/api/dynamic/%s/offline/:uuid", c.name),

			HandlerFunc: dc.offline,
		},
	}
	c.InitAccess(c.routers)
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
