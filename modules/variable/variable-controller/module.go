package variable_controller

import (
	audit_model "github.com/eolinker/apinto-dashboard/modules/audit/audit-model"
	"github.com/eolinker/apinto-module"
	"net/http"
)

type VariableDriver struct {
}

func NewVariableDriver() apinto_module.Driver {
	return &VariableDriver{}
}

func (c *VariableDriver) CreateModule(name string, config interface{}) (apinto_module.Module, error) {
	return NewVariableModule(name), nil
}

func (c *VariableDriver) CheckConfig(name string, config interface{}) error {
	return nil
}

func (c *VariableDriver) CreatePlugin(define interface{}) (apinto_module.Plugin, error) {
	return c, nil
}

type VariableModule struct {
	isInit  bool
	name    string
	routers apinto_module.RoutersInfo
}

func (c *VariableModule) Name() string {
	return c.name
}

func (c *VariableModule) Support() (apinto_module.ProviderSupport, bool) {
	return nil, false
}

func (c *VariableModule) Routers() (apinto_module.Routers, bool) {
	return c, true
}

func (c *VariableModule) Middleware() (apinto_module.Middleware, bool) {
	return nil, false
}

func NewVariableModule(name string) *VariableModule {

	return &VariableModule{name: name}
}

func (c *VariableModule) RoutersInfo() apinto_module.RoutersInfo {
	if !c.isInit {
		c.initRouter()
		c.isInit = true
	}
	return c.routers
}

func (c *VariableModule) initRouter() {
	variableController := newVariablesController()
	clusterVariableClr := newClusterVariableController()

	c.routers = []apinto_module.RouterInfo{
		{
			Method:      http.MethodGet,
			Path:        "/api/variables",
			Handler:     "variable.gets",
			HandlerFunc: []apinto_module.HandlerFunc{variableController.gets},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/variable",
			Handler:     "variable.get",
			HandlerFunc: []apinto_module.HandlerFunc{variableController.get},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/variable",
			Handler:     "variable.post",
			HandlerFunc: []apinto_module.HandlerFunc{variableController.post},
		},
		{
			Method:      http.MethodDelete,
			Path:        "/api/variable",
			Handler:     "variable.del",
			HandlerFunc: []apinto_module.HandlerFunc{variableController.del},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/cluster/:cluster_name/variables",
			Handler:     "cluster-variable.gets",
			HandlerFunc: []apinto_module.HandlerFunc{clusterVariableClr.gets},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/cluster/:cluster_name/variable",
			Handler:     "cluster-variable.post",
			HandlerFunc: []apinto_module.HandlerFunc{clusterVariableClr.post},
		},
		{
			Method:      http.MethodPut,
			Path:        "/api/cluster/:cluster_name/variable",
			Handler:     "cluster-variable.put",
			HandlerFunc: []apinto_module.HandlerFunc{clusterVariableClr.put},
		},
		{
			Method:      http.MethodDelete,
			Path:        "/api/cluster/:cluster_name/variable",
			Handler:     "cluster-variable.del",
			HandlerFunc: []apinto_module.HandlerFunc{clusterVariableClr.del},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/cluster/:cluster_name/variable/update-history",
			Handler:     "cluster-variable.updateHistory",
			HandlerFunc: []apinto_module.HandlerFunc{clusterVariableClr.updateHistory},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/cluster/:cluster_name/variable/sync-conf",
			Handler:     "cluster-variable.syncConf",
			HandlerFunc: []apinto_module.HandlerFunc{clusterVariableClr.syncConf},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/cluster/:cluster_name/variable/to-publishs",
			Handler:     "cluster-variable.toPublishs",
			HandlerFunc: []apinto_module.HandlerFunc{clusterVariableClr.toPublishs},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/cluster/:cluster_name/variable/publish",
			Handler:     "cluster-variable.publish",
			HandlerFunc: []apinto_module.HandlerFunc{audit_model.LogOperateTypePublish.Handler, clusterVariableClr.publish},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/cluster/:cluster_name/variable/publish-history",
			Handler:     "cluster-variable.publishHistory",
			HandlerFunc: []apinto_module.HandlerFunc{clusterVariableClr.publishHistory},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/cluster/:cluster_name/variable/sync-conf",
			Handler:     "cluster-variable.getSyncConf",
			HandlerFunc: []apinto_module.HandlerFunc{clusterVariableClr.getSyncConf},
		},
	}
}
