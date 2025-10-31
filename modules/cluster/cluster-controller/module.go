package cluster_controller

import (
	"github.com/eolinker/apinto-dashboard/pm3"
	"net/http"

	"github.com/eolinker/apinto-dashboard/module"
)

type ClusterPluginDriver struct {
}

func (c *ClusterPluginDriver) Install(info *pm3.PluginDefine) (ms []pm3.PModule, acs []pm3.PAccess, fs []pm3.PFrontend, err error) {

	return pm3.ReadPluginAssembly(info)

}

func (c *ClusterPluginDriver) Create(info *pm3.PluginDefine, config pm3.PluginConfig) (pm3.Module, error) {
	return NewModule(info.Id, info.Name), nil
}

func NewClusterPlugin() apinto_module.Driver {
	return &ClusterPluginDriver{}
}

type Module struct {
	*pm3.ModuleTool

	isInit bool

	name    string
	routers apinto_module.RoutersInfo
}

func (c *Module) Frontend() []pm3.FrontendAsset {
	return nil
}

func (c *Module) Apis() []pm3.Api {
	return c.RoutersInfo()
}

func (c *Module) Middleware() []pm3.Middleware {
	return nil
}

func (c *Module) Support() (pm3.ProviderSupport, bool) {
	return nil, false
}

func (c *Module) Name() string {
	return c.name
}

func NewModule(id, name string) *Module {

	return &Module{ModuleTool: pm3.NewModuleTool(id, name),
		name: name}
}

func (c *Module) RoutersInfo() apinto_module.RoutersInfo {
	if !c.isInit {
		c.initRouter()
		c.InitAccess(c.routers)
		c.isInit = true
	}
	return c.routers
}
func (c *Module) initRouter() {
	clrController := newClusterController()
	nodeController := newClusterNodeController()
	//configController := newClusterConfigController()
	certificateController := newClusterCertificateController()
	gmCertificateController := newGmCertificateController()
	c.routers = []apinto_module.RouterInfo{
		{
			Method: http.MethodGet,
			Path:   "/api/clusters",

			HandlerFunc: clrController.clusters,
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/clusters/simple",
			Authority:   pm3.Public,
			HandlerFunc: clrController.simpleClusters,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/clusters/create_check",

			HandlerFunc: clrController.createClusterCheck,
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/cluster/enum",
			Authority:   pm3.Public,
			HandlerFunc: clrController.clusterEnum,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/cluster",

			HandlerFunc: clrController.cluster,
		},
		{
			Method: http.MethodDelete,
			Path:   "/api/cluster",

			HandlerFunc: clrController.del,
		}, {
			Method: http.MethodPost,
			Path:   "/api/cluster",

			HandlerFunc: clrController.create,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/cluster-test",

			HandlerFunc: clrController.test,
		},
		{
			Method: http.MethodPut,
			Path:   "/api/cluster/:cluster_name",

			HandlerFunc: clrController.update,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/cluster/:cluster_name/nodes",

			HandlerFunc: nodeController.nodes,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/cluster/:cluster_name/nodes/simple",

			HandlerFunc: nodeController.nodesSimple,
		},
		{
			Method: http.MethodPost,
			Path:   "/api/cluster/:cluster_name/node/reset",

			HandlerFunc: nodeController.reset,
		},
		{
			Method: http.MethodPut,
			Path:   "/api/cluster/:cluster_name/node",

			HandlerFunc: nodeController.put,
		},
		//{
		//	Method:      http.MethodGet,
		//	Path:        "/api/cluster/:cluster_name/configuration/:type",
		//
		//	HandlerFunc:configController.get,
		//},
		//{
		//	Method:      http.MethodPut,
		//	Path:        "/api/cluster/:cluster_name/configuration/:type",
		//
		//	HandlerFunc:configController.edit,
		//},
		//{
		//	Method:      http.MethodPut,
		//	Path:        "/api/cluster/:cluster_name/configuration/:type/enable",
		//
		//	HandlerFunc:configController.enable,
		//},
		//{
		//	Method:      http.MethodPut,
		//	Path:        "/api/cluster/:cluster_name/configuration/:type/disable",
		//
		//	HandlerFunc:configController.disable,
		//},

		{
			Method: http.MethodPost,
			Path:   "/api/cluster/:cluster_name/certificate",

			HandlerFunc: certificateController.post,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/cluster/:cluster_name/certificate/:certificate_id",

			HandlerFunc: certificateController.get,
		},
		{
			Method: http.MethodPut,
			Path:   "/api/cluster/:cluster_name/certificate/:certificate_id",

			HandlerFunc: certificateController.put,
		},
		{
			Method: http.MethodDelete,
			Path:   "/api/cluster/:cluster_name/certificate/:certificate_id",

			HandlerFunc: certificateController.del,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/cluster/:cluster_name/certificates",

			HandlerFunc: certificateController.gets,
		},

		{
			Method: http.MethodPost,
			Path:   "/api/cluster/:cluster_name/gm_certificate",

			HandlerFunc: gmCertificateController.post,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/cluster/:cluster_name/gm_certificate/:certificate_id",

			HandlerFunc: gmCertificateController.get,
		},
		{
			Method: http.MethodPut,
			Path:   "/api/cluster/:cluster_name/gm_certificate/:certificate_id",

			HandlerFunc: gmCertificateController.put,
		},
		{
			Method: http.MethodDelete,
			Path:   "/api/cluster/:cluster_name/gm_certificate/:certificate_id",

			HandlerFunc: gmCertificateController.del,
		},
		{
			Method: http.MethodGet,
			Path:   "/api/cluster/:cluster_name/gm_certificates",

			HandlerFunc: gmCertificateController.gets,
		},
	}
}
