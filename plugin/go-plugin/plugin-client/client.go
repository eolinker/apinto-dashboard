package plugin_client

import (
	"context"
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/plugin/go-plugin/proto"
	"github.com/eolinker/apinto-dashboard/plugin/go-plugin/shared"
	"github.com/eolinker/apinto-dashboard/pm3"
	"github.com/eolinker/eosc/log"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

var _ ClientHandler = (*GrpcClient)(nil)

type ClientHandler = pm3.Module

type GrpcClient struct {
	*pm3.ModuleTool

	name       string
	client     proto.ServiceClient
	apis       []pm3.Api
	middleware []pm3.Middleware
	frontends  []pm3.FrontendAsset
}

func (g *GrpcClient) init(ctx context.Context) error {
	info, err := g.client.ModuleInfo(ctx, &proto.Empty{})
	if err != nil {
		return err
	}
	g.apis = make([]pm3.Api, 0, len(info.Apis))
	for _, a := range info.Apis {
		g.apis = append(g.apis, g.genApi(a))
	}
	g.InitAccess(g.apis)

	g.middleware = make([]pm3.Middleware, 0, len(info.Middlewares))
	for _, m := range info.Middlewares {
		g.middleware = append(g.middleware, g.genMiddlewareHandler(m))
	}

	g.frontends = make([]pm3.FrontendAsset, 0, len(info.Frontend))

	for _, f := range info.Frontend {
		p := newProxy(g.client, f.Id)
		g.frontends = append(g.frontends, pm3.FrontendAsset{
			Path:        f.Path,
			HandlerFunc: p.Handle,
		})
	}
	return nil
}

type MiddlewareHandler struct {
	hasRequestHandler  bool
	hasResponseHandler bool
	name               string
	client             proto.ServiceClient
}

func (g *MiddlewareHandler) Handle(ginCtx *gin.Context) {
	if g.hasRequestHandler {
		req := &proto.MiddlewareRequest{
			Name:    g.name,
			Request: readRequest(ginCtx),
		}

		response, err := g.client.MiddlewaresRequest(ginCtx, req)
		if err != nil {

			return
		}

		if !middlewareDoResponse(response, ginCtx) {
			return
		}
	}
	if g.hasResponseHandler {
		ginCtx.Next()

		req := &proto.MiddlewareRequest{
			Name:    g.name,
			Request: readRequest(ginCtx),
		}

		mRespResponse, err := g.client.MiddlewaresResponse(ginCtx, req)
		if err != nil {
			return
		}

		middlewareDoResponse(mRespResponse, ginCtx)
	}

}

func (m *MiddlewareHandler) Check(api pm3.ApiInfo) bool {
	r, err := m.client.CheckMiddlewareForApi(context.Background(), &proto.MiddlewareInfoRequest{
		Name: m.name,
		Api: &proto.ApiInfo{
			Method:    api.Method,
			Path:      api.Path,
			Access:    api.Access,
			Authority: proto.Authority(api.Authority),
			Id:        "",
		},
	})
	if err != nil {
		return false
	}
	return r.Result
}

func (g *GrpcClient) genMiddlewareHandler(info *proto.MiddlewareInfo) *MiddlewareHandler {
	return &MiddlewareHandler{name: info.Name, client: g.client, hasRequestHandler: info.HasRequestHandler, hasResponseHandler: info.HasResponseHandler}
}

func (g *GrpcClient) genApi(info *proto.ApiInfo) pm3.Api {
	p := newProxy(g.client, info.Id)
	return pm3.Api{

		Authority: pm3.ApiAuthority(info.Authority),
		Access:    info.Access,
		Method:    info.Method,
		Path:      info.Path,

		HandlerFunc: p.Handle,
	}
}
func (g *GrpcClient) Name() string {
	return g.name
}

func (g *GrpcClient) Frontend() []pm3.FrontendAsset {
	return g.frontends
}

func (g *GrpcClient) Apis() []pm3.Api {
	return g.apis
}

func (g *GrpcClient) Middleware() []pm3.Middleware {
	return g.middleware
}

func (g *GrpcClient) Support() (pm3.ProviderSupport, bool) {
	return nil, false
}

func readRequest(ginCtx *gin.Context) *proto.HttpRequest {
	r := &proto.HttpRequest{
		Url:     ginCtx.Request.RequestURI,
		Method:  ginCtx.Request.Method,
		FulPath: ginCtx.FullPath(),
	}
	r.Body, _ = ginCtx.GetRawData()
	ginCtx.Request.Body.Close()
	if ginCtx.Keys != nil {
		r.Keys, _ = json.Marshal(ginCtx.Keys)
	}
	r.Headers = make([]*proto.Header, 0, len(ginCtx.Request.Header))
	for k, v := range ginCtx.Request.Header {
		r.Headers = append(r.Headers, &proto.Header{
			Key:   k,
			Value: v,
		})
	}
	return r
}

func middlewareDoResponse(response *proto.MiddlewareResponse, ginCtx *gin.Context) bool {
	if len(response.Headers) > 0 {
		for _, h := range response.Headers {
			if len(h.Value) == 1 {
				ginCtx.Header(h.Key, h.Value[0])
			} else {
				ginCtx.Writer.Header().Del(h.Key)
				for _, v := range h.Value {
					ginCtx.Writer.Header().Add(h.Key, v)
				}
			}

		}
	}

	if len(response.Keys) > 0 {
		keys := make(map[string]interface{})
		err := json.Unmarshal(response.Keys, &keys)
		if err != nil {
			log.Warn(err.Error())
		} else {
			for k, v := range keys {
				ginCtx.Set(k, v)
			}
		}
	}
	if response.Abort {
		ginCtx.Abort()
		switch response.Action {
		case "redirect":
			ginCtx.Redirect(int(response.StatusCode), string(response.Body))
			return false
		default:

			ginCtx.Data(int(response.StatusCode), response.ContentType, response.Body)
			return true
		}
	}
	return true
}

type GrpcPluginClient struct {
	plugin.Plugin
	name string
	id   string
}

func (g *GrpcPluginClient) GRPCServer(broker *plugin.GRPCBroker, server *grpc.Server) error {
	//TODO implement me
	panic("implement me")
}

func (g *GrpcPluginClient) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, conn *grpc.ClientConn) (interface{}, error) {
	client := &GrpcClient{
		ModuleTool: pm3.NewModuleTool(g.id, g.name),
		client:     proto.NewServiceClient(conn), name: g.name}
	err := client.init(ctx)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func CreateClientProxy(id, name string) map[string]plugin.Plugin {
	return map[string]plugin.Plugin{

		shared.PluginHandlerName: &GrpcPluginClient{
			id:   id,
			name: name,
		},
	}
}
