package plugin

import (
	apinto_module "github.com/eolinker/apinto-dashboard/module"
	proto2 "github.com/eolinker/apinto-dashboard/plugin/go-plugin/proto"
	"github.com/eolinker/apinto-dashboard/plugin/go-plugin/shared"
	"github.com/eolinker/eosc/log"
	"github.com/gin-gonic/gin"
	"github.com/go-basic/uuid"
	"github.com/hashicorp/go-plugin"
)

type Plugin struct {
	//middlewares map[string]*serverMiddlewareHandler
	module shared.Module
	server *Server
}

func NewPlugin(module shared.Module) *Plugin {

	return &Plugin{module: module}
}

func (b *Plugin) Server() {
	if b.server == nil {
		err := b.build()
		if err != nil {
			panic(err)
		}

	}
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.HandshakeConfig,
		TLSProvider:     nil,
		Plugins: map[string]plugin.Plugin{
			shared.PluginHandlerName: &GrpcPluginServer{
				iml: b.server,
			},
		},
		VersionedPlugins: nil,
		GRPCServer:       plugin.DefaultGRPCServer,
		Logger:           logger,
		Test:             nil,
	})
}

func (b *Plugin) build() error {

	infos := &proto2.PluginInfos{}
	handlers := make(map[string]gin.HandlerFunc)
	engine := Engine()
	apis := b.module.Apis()
	infos.Apis = make([]*proto2.ApiInfo, 0, len(apis))
	for _, a := range b.module.Apis() {
		id := uuid.New()
		engine.Handle(a.Method, a.Path, a.HandlerFunc)
		infos.Apis = append(infos.Apis, &proto2.ApiInfo{
			Method:    a.Method,
			Path:      a.Path,
			Access:    a.Access,
			Authority: proto2.Authority(a.Authority),
			Id:        id,
		})

		handlers[id] = a.HandlerFunc
	}
	frontends := b.module.Frontend()
	infos.Frontend = make([]*proto2.FrontendInfo, 0, len(frontends))
	for _, f := range frontends {
		path := apinto_module.StaticRouter(f.Path)
		engine.GET(path, f.HandlerFunc)
		engine.HEAD(path, f.HandlerFunc)

		id := uuid.New()
		infos.Frontend = append(infos.Frontend, &proto2.FrontendInfo{
			Path: f.Path,
			Id:   id,
		})
		handlers[id] = f.HandlerFunc

	}
	infos.Middlewares = make([]*proto2.MiddlewareInfo, 0, len(b.module.Middleware()))
	middlewares := make(map[string]*serverMiddlewareHandler)
	for _, m := range b.module.Middleware() {
		name := uuid.New()
		requestH, hasRequest := m.RequestHandler()
		responseH, hasResponse := m.ResponseHandler()

		infos.Middlewares = append(infos.Middlewares, &proto2.MiddlewareInfo{
			Name:               name,
			HasRequestHandler:  hasRequest,
			HasResponseHandler: hasResponse,
		})
		middlewares[name] = &serverMiddlewareHandler{
			requestHandler:  requestH,
			responseHandler: responseH,
			checkHandler:    m.Check,
		}
	}
	for _, r := range engine.Routes() {
		log.DebugF("router:%s %s %s", r.Method, r.Path, r.Handler)
	}

	b.server = &Server{
		httpHandler: engine,
		middlewares: middlewares,
		infos:       infos,
		handlers:    handlers,
	}
	return nil

}
