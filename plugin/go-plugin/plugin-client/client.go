package plugin_client

import (
	"context"
	"encoding/json"
	proto "github.com/eolinker/apinto-dashboard/plugin/go-plugin/proto"
	"github.com/eolinker/apinto-dashboard/plugin/go-plugin/shared"
	"github.com/eolinker/eosc/log"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

var _ ClientHandler = (*GrpcClient)(nil)

type ResponseHandler func(ctx context.Context, response *proto.HttpResponse)
type RequestHandler func(ctx context.Context, request *proto.HttpRequest)

type MiddlewareResponseHandler func(ctx context.Context, response *proto.MiddlewareResponse)
type MiddlewareRequestHandler func(ctx context.Context, request *proto.MiddlewareRequest)

type ClientHandler interface {
	ServerGin(ginCtx *gin.Context, requestHandler RequestHandler, responseHandler ResponseHandler)
	CreateMiddleware(name string, requestHandler MiddlewareRequestHandler, responseHandler MiddlewareResponseHandler) (*MiddlewareClientHandler, error)
}

// PluginMap is the map of plugins we can dispense.
var PluginMap = map[string]plugin.Plugin{}

type MiddlewareClientHandler struct {
	client          proto.ServiceClient
	name            string
	request         bool
	response        bool
	requestHandler  MiddlewareRequestHandler
	responseHandler MiddlewareResponseHandler
}
type GrpcClient struct {
	client proto.ServiceClient
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
func (g *GrpcClient) ServerGin(ginCtx *gin.Context, requestHandler RequestHandler, responseHandler ResponseHandler) {
	req := readRequest(ginCtx)
	if requestHandler != nil {
		requestHandler(ginCtx, req)
	}
	resp, err := g.client.Request(ginCtx, req)
	if err != nil {
		log.Errorf(ginCtx.AbortWithError(500, err).Error())
		return
	}
	if responseHandler != nil {
		responseHandler(ginCtx, resp)
	}
	header := ginCtx.Writer.Header()
	for _, h := range resp.Headers {
		for _, v := range h.Value {
			header.Add(h.Key, v)
		}
	}

	contentType := header.Get("content-type")
	ginCtx.Data(int(resp.Status), contentType, resp.Body)
}
func (g *GrpcClient) CreateMiddleware(name string, requestHandler MiddlewareRequestHandler, responseHandler MiddlewareResponseHandler) (*MiddlewareClientHandler, error) {
	info, err := g.client.GetMiddlewareInfo(context.Background(), &proto.MiddlewareInfoRequest{Name: name})
	if err != nil {
		return nil, err
	}
	return &MiddlewareClientHandler{name: name, client: g.client, request: info.Request, response: info.Response, requestHandler: requestHandler, responseHandler: responseHandler}, nil
}
func (g *MiddlewareClientHandler) Middleware(ginCtx *gin.Context) {
	if g.request {
		req := &proto.MiddlewareRequest{
			Name:    g.name,
			Request: readRequest(ginCtx),
		}
		if g.requestHandler != nil {
			g.requestHandler(ginCtx, req)
		}
		response, err := g.client.MiddlewaresRequest(ginCtx, req)
		if err != nil {

			return
		}
		if g.responseHandler != nil {
			g.responseHandler(ginCtx, response)
		}
		if !middlewareDoResponse(response, ginCtx) {
			return
		}
	}
	if g.response {
		ginCtx.Next()

		req := &proto.MiddlewareRequest{
			Name:    g.name,
			Request: readRequest(ginCtx),
		}
		if g.requestHandler != nil {
			g.requestHandler(ginCtx, req)
		}
		mRespResponse, err := g.client.MiddlewaresResponse(ginCtx, req)
		if err != nil {
			return
		}
		if g.responseHandler != nil {
			g.responseHandler(ginCtx, mRespResponse)
		}
		middlewareDoResponse(mRespResponse, ginCtx)
	}

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
}

func (g *GrpcPluginClient) GRPCServer(broker *plugin.GRPCBroker, server *grpc.Server) error {
	//TODO implement me
	panic("implement me")
}

func (g *GrpcPluginClient) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, conn *grpc.ClientConn) (interface{}, error) {
	return &GrpcClient{client: proto.NewServiceClient(conn)}, nil
}

func CreateClientProxy() map[string]plugin.Plugin {
	return map[string]plugin.Plugin{
		shared.PluginHandlerName: &GrpcPluginClient{},
	}
}
