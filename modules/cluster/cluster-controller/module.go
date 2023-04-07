package cluster_controller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/eolinker/apinto-module"
)

type ClusterPluginDriver struct {
}

func NewClusterPlugin() apinto_module.Driver {
	return &ClusterPluginDriver{}
}

const (
	ClusterView = "view"
	ClusterEdit = "edit"
)

func (c *ClusterPluginDriver) CreateModule(name string, apiPrefix string, config interface{}) (apinto_module.Module, error) {
	return NewClusterModule(apiPrefix), nil
}

func (c *ClusterPluginDriver) CheckConfig(name string, apiPrefix string, config interface{}) error {
	return nil
}

func (c *ClusterPluginDriver) CreatePlugin(define interface{}) (apinto_module.Plugin, error) {
	return c, nil
}

type ClusterModule struct {
	isInit    bool
	apiPrefix string
	name      string
	routers   apinto_module.RoutersInfo
}

func (c *ClusterModule) Access() []apinto_module.AccessInfo {
	//TODO implement me
	panic("implement me")
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

func NewClusterModule(apiPrefix string) *ClusterModule {
	if !strings.HasPrefix(apiPrefix, "/") {
		apiPrefix = "/" + apiPrefix
	}
	apiPrefix = strings.TrimSuffix(apiPrefix, "/")

	return &ClusterModule{apiPrefix: apiPrefix}
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
			Path:        fmt.Sprintf("%s/clusters", c.apiPrefix),
			Handler:     "cluster.list",
			HandlerFunc: clrController.clusters,
		},
		{
			Method:      http.MethodGet,
			Path:        fmt.Sprintf("%s/cluster/enum", c.apiPrefix),
			Handler:     "cluster.enum",
			HandlerFunc: clrController.clusterEnum,
		},
		{
			Method:      http.MethodGet,
			Path:        fmt.Sprintf("%s/cluster", c.apiPrefix),
			Handler:     "cluster.info",
			HandlerFunc: clrController.cluster,
		},
		{
			Method:      http.MethodDelete,
			Path:        fmt.Sprintf("%s/cluster", c.apiPrefix),
			Handler:     "cluster.delete",
			HandlerFunc: clrController.del,
		}, {
			Method:      http.MethodPost,
			Path:        fmt.Sprintf("%s/cluster", c.apiPrefix),
			Handler:     "cluster.create",
			HandlerFunc: clrController.create,
		},
		{
			Method:      http.MethodGet,
			Path:        fmt.Sprintf("%s/cluster-test", c.apiPrefix),
			Handler:     "cluster.test",
			HandlerFunc: clrController.test,
		},
		{
			Method:      http.MethodPut,
			Path:        fmt.Sprintf("%s/cluster/:cluster_name/desc", c.apiPrefix),
			Handler:     "cluster.desc.edit",
			HandlerFunc: clrController.putDesc,
		},
		{
			Method:      http.MethodGet,
			Path:        fmt.Sprintf("%s/cluster/:cluster_name/nodes", c.apiPrefix),
			Handler:     "cluster.nodes",
			HandlerFunc: nodeController.nodes,
		},
		{
			Method:      http.MethodPost,
			Path:        fmt.Sprintf("%s/cluster/:cluster_name/node/reset", c.apiPrefix),
			Handler:     "cluster.nodes.reset",
			HandlerFunc: nodeController.reset,
		},
		{
			Method:      http.MethodPut,
			Path:        fmt.Sprintf("%s/cluster/:cluster_name/node", c.apiPrefix),
			Handler:     "cluster.nodes.edit",
			HandlerFunc: nodeController.put,
		},
		{
			Method:      http.MethodGet,
			Path:        fmt.Sprintf("%s/cluster/:cluster_name/configuration/:type", c.apiPrefix),
			Handler:     "cluster.config",
			HandlerFunc: configController.get,
		},
		{
			Method:      http.MethodPut,
			Path:        fmt.Sprintf("%s/cluster/:cluster_name/configuration/:type", c.apiPrefix),
			Handler:     "cluster.config.edit",
			HandlerFunc: configController.edit,
		},
		{
			Method:      http.MethodPut,
			Path:        fmt.Sprintf("%s/cluster/:cluster_name/configuration/:type/enable", c.apiPrefix),
			Handler:     "cluster.config.enable",
			HandlerFunc: configController.enable,
		},
		{
			Method:      http.MethodPut,
			Path:        fmt.Sprintf("%s/cluster/:cluster_name/configuration/:type/disable", c.apiPrefix),
			Handler:     "cluster.config.disable",
			HandlerFunc: configController.disable,
		},
	}
}
