package variable_controller

import (
	"github.com/eolinker/apinto-dashboard/module"
	"github.com/eolinker/apinto-dashboard/pm3"
	"net/http"
)

type VariableDriver struct {
}

func (c *VariableDriver) Install(info *pm3.PluginDefine) (ms []pm3.PModule, acs []pm3.PAccess, fs []pm3.PFrontend, err error) {
	return pm3.ReadPluginAssembly(info)
}

func (c *VariableDriver) Create(info *pm3.PluginDefine, config pm3.PluginConfig) (pm3.Module, error) {
	return NewVariableModule(info.Id, info.Name), nil
}

func NewVariableDriver() apinto_module.Driver {
	return &VariableDriver{}
}

//
//
//func (c *VariableDriver) GetPluginFrontend(moduleName string) string {
//	return "deploy/variable"
//}

type VariableModule struct {
	*pm3.ModuleTool

	isInit  bool
	name    string
	routers apinto_module.RoutersInfo
}

func (c *VariableModule) Frontend() []pm3.FrontendAsset {
	return nil
}

func (c *VariableModule) Apis() []pm3.Api {
	if !c.isInit {
		c.initRouter()
		c.InitAccess(c.routers)

		c.isInit = true
	}
	return c.routers
}

func (c *VariableModule) Middleware() []pm3.Middleware {
	return nil
}

func (c *VariableModule) Name() string {
	return c.name
}

func (c *VariableModule) Support() (pm3.ProviderSupport, bool) {
	return nil, false
}

func NewVariableModule(id, name string) *VariableModule {

	return &VariableModule{name: name, ModuleTool: pm3.NewModuleTool(id, name)}
}

func (c *VariableModule) initRouter() {
	variableController := newVariablesController()
	clusterVariableClr := newClusterVariableController()

	c.routers = []apinto_module.RouterInfo{
		{
			Method: http.MethodGet,
			Path:   "/api/variables",

			HandlerFunc: variableController.gets,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/variable",

			HandlerFunc: variableController.get,
		},
		{
			Method: http.MethodPost,
			Path:   "/api/variable",

			HandlerFunc: variableController.post,
		},
		{
			Method: http.MethodDelete,
			Path:   "/api/variable",

			HandlerFunc: variableController.del,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/cluster/:cluster_name/variables",

			HandlerFunc: clusterVariableClr.gets,
		},
		{
			Method: http.MethodPost,
			Path:   "/api/cluster/:cluster_name/variable",

			HandlerFunc: clusterVariableClr.post,
		},
		{
			Method: http.MethodPut,
			Path:   "/api/cluster/:cluster_name/variable",

			HandlerFunc: clusterVariableClr.put,
		},
		{
			Method: http.MethodDelete,
			Path:   "/api/cluster/:cluster_name/variable",

			HandlerFunc: clusterVariableClr.del,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/cluster/:cluster_name/variable/update-history",

			HandlerFunc: clusterVariableClr.updateHistory,
		},
		{
			Method: http.MethodPost,
			Path:   "/api/cluster/:cluster_name/variable/sync-conf",

			HandlerFunc: clusterVariableClr.syncConf,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/cluster/:cluster_name/variable/to-publishs",

			HandlerFunc: clusterVariableClr.toPublishs,
		},
		{
			Method: http.MethodPost,
			Path:   "/api/cluster/:cluster_name/variable/publish",

			HandlerFunc: clusterVariableClr.publish,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/cluster/:cluster_name/variable/publish-history",

			HandlerFunc: clusterVariableClr.publishHistory,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/cluster/:cluster_name/variable/sync-conf",

			HandlerFunc: clusterVariableClr.getSyncConf,
		},
	}
}
