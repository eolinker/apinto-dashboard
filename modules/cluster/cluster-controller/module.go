package cluster_controller

import (
	"net/http"

	"github.com/eolinker/apinto-dashboard/module"
)

type ClusterPluginDriver struct {
}

func NewClusterPlugin() apinto_module.Driver {
	return &ClusterPluginDriver{}
}

func (c *ClusterPluginDriver) CreateModule(name string, config interface{}) (apinto_module.Module, error) {
	return NewModule(name), nil
}

func (c *ClusterPluginDriver) CheckConfig(name string, config interface{}) error {
	return nil
}

func (c *ClusterPluginDriver) CreatePlugin(define interface{}) (apinto_module.Plugin, error) {
	return c, nil
}

func (c *ClusterPluginDriver) GetPluginFrontend(moduleName string) string {
	return "deploy/cluster"
}

func (c *ClusterPluginDriver) IsPluginVisible() bool {
	return true
}

func (c *ClusterPluginDriver) IsShowServer() bool {
	return false
}

func (c *ClusterPluginDriver) IsCanUninstall() bool {
	return false
}

func (c *ClusterPluginDriver) IsCanDisable() bool {
	return false
}

type Module struct {
	isInit bool

	name    string
	routers apinto_module.RoutersInfo
}

func (c *Module) Name() string {
	return c.name
}

func (c *Module) Support() (apinto_module.ProviderSupport, bool) {
	return nil, false
}

func (c *Module) Routers() (apinto_module.Routers, bool) {
	return c, true
}

func (c *Module) Middleware() (apinto_module.Middleware, bool) {
	return nil, false
}

func NewModule(name string) *Module {

	return &Module{name: name}
}

func (c *Module) RoutersInfo() apinto_module.RoutersInfo {
	if !c.isInit {
		c.initRouter()
		c.isInit = true
	}
	return c.routers
}
func (c *Module) initRouter() {
	clrController := newClusterController()
	nodeController := newClusterNodeController()
	//configController := newClusterConfigController()
	certificateController := newClusterCertificateController()
	c.routers = []apinto_module.RouterInfo{
		{
			Method:      http.MethodGet,
			Path:        "/api/clusters",
			Handler:     "cluster.list",
			HandlerFunc: []apinto_module.HandlerFunc{clrController.clusters},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/clusters/simple",
			Handler:     "cluster.simple_list",
			HandlerFunc: []apinto_module.HandlerFunc{clrController.simpleClusters},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/clusters/create_check",
			Handler:     "cluster.simple_list",
			HandlerFunc: []apinto_module.HandlerFunc{clrController.createClusterCheck},
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
			Path:        "/api/cluster/:cluster_name",
			Handler:     "cluster.desc.edit",
			HandlerFunc: []apinto_module.HandlerFunc{clrController.update},
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
		//{
		//	Method:      http.MethodGet,
		//	Path:        "/api/cluster/:cluster_name/configuration/:type",
		//	Handler:     "cluster.config",
		//	HandlerFunc: []apinto_module.HandlerFunc{configController.get},
		//},
		//{
		//	Method:      http.MethodPut,
		//	Path:        "/api/cluster/:cluster_name/configuration/:type",
		//	Handler:     "cluster.config.edit",
		//	HandlerFunc: []apinto_module.HandlerFunc{configController.edit},
		//},
		//{
		//	Method:      http.MethodPut,
		//	Path:        "/api/cluster/:cluster_name/configuration/:type/enable",
		//	Handler:     "cluster.config.enable",
		//	HandlerFunc: []apinto_module.HandlerFunc{configController.enable},
		//},
		//{
		//	Method:      http.MethodPut,
		//	Path:        "/api/cluster/:cluster_name/configuration/:type/disable",
		//	Handler:     "cluster.config.disable",
		//	HandlerFunc: []apinto_module.HandlerFunc{configController.disable},
		//},

		{
			Method:      http.MethodPost,
			Path:        "/api/cluster/:cluster_name/certificate",
			Handler:     "cluster.certificates.post",
			HandlerFunc: []apinto_module.HandlerFunc{certificateController.post},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/cluster/:cluster_name/certificate/:certificate_id",
			Handler:     "cluster.certificates.get",
			HandlerFunc: []apinto_module.HandlerFunc{certificateController.get},
		},
		{
			Method:      http.MethodPut,
			Path:        "/api/cluster/:cluster_name/certificate/:certificate_id",
			Handler:     "cluster.certificates.put",
			HandlerFunc: []apinto_module.HandlerFunc{certificateController.put},
		},
		{
			Method:      http.MethodDelete,
			Path:        "/api/cluster/:cluster_name/certificate/:certificate_id",
			Handler:     "cluster.certificates.del",
			HandlerFunc: []apinto_module.HandlerFunc{certificateController.del},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/cluster/:cluster_name/certificates",
			Handler:     "cluster.certificates.gets",
			HandlerFunc: []apinto_module.HandlerFunc{certificateController.gets},
		},
	}
}
