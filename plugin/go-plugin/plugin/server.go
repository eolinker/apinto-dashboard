package plugin

import (
	"bytes"
	"context"
	"encoding/json"
	module "github.com/eolinker/apinto-dashboard/module"
	proto2 "github.com/eolinker/apinto-dashboard/plugin/go-plugin/proto"
	"github.com/eolinker/apinto-dashboard/plugin/go-plugin/shared"
	"github.com/eolinker/apinto-dashboard/pm3"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-plugin"
	"io"
	"net/http"
	"net/url"

	"google.golang.org/grpc"
)

type serverMiddlewareHandler struct {
	requestHandler  shared.MiddlewareHandFunc
	responseHandler shared.MiddlewareHandFunc
	checkHandler    func(api pm3.ApiInfo) bool
}

var (
	_ proto2.ServiceServer = (*Server)(nil)
)

type Server struct {
	proto2.UnimplementedServiceServer
	httpHandler http.Handler
	middlewares map[string]*serverMiddlewareHandler

	infos    *proto2.PluginInfos
	handlers map[string]gin.HandlerFunc
}

func (p *Server) ModuleInfo(ctx context.Context, empty *proto2.Empty) (*proto2.PluginInfos, error) {
	return p.infos, nil
}

func (p *Server) CheckMiddlewareForApi(ctx context.Context, request *proto2.MiddlewareInfoRequest) (*proto2.MiddlewareInfoResponse, error) {
	result := false
	if h, has := p.middlewares[request.Name]; has {
		result = h.checkHandler(pm3.ApiInfo{
			Authority: pm3.ApiAuthority(request.Api.Authority),
			Access:    request.Api.Access,
			Method:    request.Api.Method,
			Path:      request.Api.Path,
		})

	}
	return &proto2.MiddlewareInfoResponse{Result: result}, nil
}

func (p *Server) MiddlewaresRequest(ctx context.Context, request *proto2.MiddlewareRequest) (*proto2.MiddlewareResponse, error) {

	m, ok := p.middlewares[request.Name]
	if !ok {
		return new(proto2.MiddlewareResponse), nil
	}

	return doMiddlewareHandler(m.requestHandler, ctx, request)
}

func (p *Server) MiddlewaresResponse(ctx context.Context, request *proto2.MiddlewareRequest) (*proto2.MiddlewareResponse, error) {
	response := new(proto2.MiddlewareResponse)
	m, ok := p.middlewares[request.Name]
	if !ok {
		return response, nil
	}

	return doMiddlewareHandler(m.responseHandler, ctx, request)

}
func doMiddlewareHandler(handFunc shared.MiddlewareHandFunc, ctx context.Context, request *proto2.MiddlewareRequest) (*proto2.MiddlewareResponse, error) {
	response := new(proto2.MiddlewareResponse)
	if handFunc == nil {
		return response, nil
	}
	r := &module.MiddlewareRequest{
		FulPath: request.Request.FulPath,
		Url:     request.Request.Url,
		Method:  request.Request.Method,
		Header:  make(map[string][]string),
		Keys:    make(map[string]any),
	}

	for _, h := range request.Request.Headers {
		r.Header[h.Key] = h.Value
	}
	if len(request.Request.Keys) > 0 {
		json.Unmarshal(request.Request.Keys, &r.Keys)
	}

	writer := new(module.MiddlewareResponse)
	handFunc(ctx, r, writer)

	response.StatusCode = int32(writer.StatusCode)
	response.Body = writer.Body
	response.Headers = make([]*proto2.Header, 0, len(writer.Header))
	for k, v := range writer.Header {
		response.Headers = append(response.Headers, &proto2.Header{Key: k, Value: v})
	}
	response.Action = writer.Action
	response.Abort = writer.Abort
	if writer.Keys != nil {
		response.Keys, _ = json.Marshal(writer.Keys)
	}
	response.ContentType = writer.ContentType

	return response, nil
}
func (p *Server) GetMiddlewareInfo(ctx context.Context, r *proto2.MiddlewareInfoRequest) (*proto2.MiddlewareInfoResponse, error) {
	info := &proto2.MiddlewareInfoResponse{
		Result: false,
	}
	if r.GetApi() == nil {
		return info, nil
	}
	m, ok := p.middlewares[r.Name]
	if !ok {
		return info, nil
	}
	m.checkHandler(pm3.ApiInfo{
		Authority: pm3.ApiAuthority(r.Api.Authority),
		Access:    r.Api.Access,
		Method:    r.Api.Method,
		Path:      r.Api.Path,
	})
	return info, nil
}

func (p *Server) Request(ctx context.Context, request *proto2.HttpRequest) (*proto2.HttpResponse, error) {
	uri, err := url.ParseRequestURI(request.Url)
	if err != nil {
		return nil, err
	}
	if len(request.Keys) > 0 {
		keys := make(map[string]interface{})
		err := json.Unmarshal(request.Keys, &keys)
		if err == nil {
			ctx = widthKey(ctx, keys)
		}
	}
	req := &http.Request{
		Method: request.Method,
		URL:    uri,

		Header: make(http.Header),
	}
	req = req.WithContext(ctx)
	for _, h := range request.Headers {
		req.Header[h.Key] = h.Value
	}

	req.Body = io.NopCloser(bytes.NewReader(request.Body))

	w := newResponseWriter()

	p.httpHandler.ServeHTTP(w, req)

	responseHeader := make([]*proto2.Header, len(w.Header()))
	for k, v := range w.header {
		responseHeader = append(responseHeader, &proto2.Header{Key: k, Value: v})
	}
	return &proto2.HttpResponse{
		Status:  int32(w.status),
		Body:    w.buf.Bytes(),
		Headers: responseHeader,
	}, nil
}

type GrpcPluginServer struct {
	plugin.Plugin
	iml *Server
}

func (p *GrpcPluginServer) GRPCServer(broker *plugin.GRPCBroker, server *grpc.Server) error {
	proto2.RegisterServiceServer(server, p.iml)
	return nil
}

func (p *GrpcPluginServer) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, conn *grpc.ClientConn) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}
