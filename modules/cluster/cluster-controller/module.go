package cluster_controller

import (
	"github.com/eolinker/apinto-module"
	"net/http"
)

type ClusterPluginDriver struct {
}

func NewClusterPlugin() apinto_module.Driver {
	return &ClusterPluginDriver{}
}

func (c *ClusterPluginDriver) CreateModule(name string, config interface{}) (apinto_module.Module, error) {
	return NewClusterModule(name), nil
}

func (c *ClusterPluginDriver) CheckConfig(name string, config interface{}) error {
	return nil
}

func (c *ClusterPluginDriver) CreatePlugin(define interface{}) (apinto_module.Plugin, error) {
	return c, nil
}

type ClusterModule struct {
	isInit bool

	name    string
	routers apinto_module.RoutersInfo
}

func (c *ClusterModule) Name() string {
	return c.name
}

func (c *ClusterModule) Support() (apinto_module.ProviderSupport, bool) {
	return nil, false
}

func (c *ClusterModule) Routers() (apinto_module.Routers, bool) {
	return c, true
}

func (c *ClusterModule) Middleware() (apinto_module.Middleware, bool) {
	return nil, false
}

func NewClusterModule(name string) *ClusterModule {

	return &ClusterModule{name: name}
}

func (c *ClusterModule) RoutersInfo() apinto_module.RoutersInfo {
	if !c.isInit {
		c.initRouter()
		c.isInit = true
	}
	return c.routers
}
func (c *ClusterModule) initRouter() {
	clrController := newClusterController()
	nodeController := newClusterNodeController()
	configController := newClusterConfigController()

	c.routers = []apinto_module.RouterInfo{
		{
			Method:      http.MethodGet,
			Path:        "/api/cluster/lists",
			Handler:     "cluster.list",
			HandlerFunc: []apinto_module.HandlerFunc{clrController.clusters},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/cluster/enum",
			Handler:     "cluster.enum",
			HandlerFunc: []apinto_module.HandlerFunc{clrController.clusterEnum},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/cluster",
			Handler:     "cluster.info",
			HandlerFunc: []apinto_module.HandlerFunc{clrController.cluster},
		},
		{
			Method:      http.MethodDelete,
			Path:        "/api/cluster",
			Handler:     "cluster.delete",
			HandlerFunc: []apinto_module.HandlerFunc{clrController.del},
		}, {
			Method:      http.MethodPost,
			Path:        "/api/cluster",
			Handler:     "cluster.create",
			HandlerFunc: []apinto_module.HandlerFunc{clrController.create},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/cluster-test",
			Handler:     "cluster.test",
			HandlerFunc: []apinto_module.HandlerFunc{clrController.test},
		},
		{
			Method:      http.MethodPut,
			Path:        "/api/cluster/:cluster_name/desc",
			Handler:     "cluster.desc.edit",
			HandlerFunc: []apinto_module.HandlerFunc{clrController.putDesc},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/cluster/:cluster_name/nodes",
			Handler:     "cluster.nodes",
			HandlerFunc: []apinto_module.HandlerFunc{nodeController.nodes},
		},
		{
			Method:      http.MethodPost,
			Path:        "/api/cluster/:cluster_name/node/reset",
			Handler:     "cluster.nodes.reset",
			HandlerFunc: []apinto_module.HandlerFunc{nodeController.reset},
		},
		{
			Method:      http.MethodPut,
			Path:        "/api/cluster/:cluster_name/node",
			Handler:     "cluster.nodes.edit",
			HandlerFunc: []apinto_module.HandlerFunc{nodeController.put},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/cluster/:cluster_name/configuration/:type",
			Handler:     "cluster.config",
			HandlerFunc: []apinto_module.HandlerFunc{configController.get},
		},
		{
			Method:      http.MethodPut,
			Path:        "/api/cluster/:cluster_name/configuration/:type",
			Handler:     "cluster.config.edit",
			HandlerFunc: []apinto_module.HandlerFunc{configController.edit},
		},
		{
			Method:      http.MethodPut,
			Path:        "/api/cluster/:cluster_name/configuration/:type/enable",
			Handler:     "cluster.config.enable",
			HandlerFunc: []apinto_module.HandlerFunc{configController.enable},
		},
		{
			Method:      http.MethodPut,
			Path:        "/api/cluster/:cluster_name/configuration/:type/disable",
			Handler:     "cluster.config.disable",
			HandlerFunc: []apinto_module.HandlerFunc{configController.disable},
		},
	}
}
