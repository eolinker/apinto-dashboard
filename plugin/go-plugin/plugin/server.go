package plugin

import (
	"bytes"
	"context"
	"encoding/json"
	module "github.com/eolinker/apinto-dashboard/module"
	proto2 "github.com/eolinker/apinto-dashboard/plugin/go-plugin/proto"
	"github.com/eolinker/apinto-dashboard/plugin/go-plugin/shared"
	"github.com/hashicorp/go-plugin"
	"io"
	"net/http"
	"net/url"

	"google.golang.org/grpc"
)

type MiddlewareHandFunc func(ctx context.Context, request *module.MiddlewareRequest, writer module.MiddlewareResponseWriter)

type serverMiddlewareHandler struct {
	requestHandler  MiddlewareHandFunc
	responseHandler MiddlewareHandFunc
}

type Plugin struct {
	httpHandler shared.HttpHandler
	middlewares map[string]*serverMiddlewareHandler
}

func NewPlugin(httpHandler shared.HttpHandler) *Plugin {
	return &Plugin{httpHandler: httpHandler, middlewares: make(map[string]*serverMiddlewareHandler)}
}

type SetFunc func(builder *serverMiddlewareHandler)

func (b *Plugin) AddMiddleware(name string, setFuncs ...SetFunc) {
	m := &serverMiddlewareHandler{}
	for _, setFunc := range setFuncs {
		setFunc(m)
	}
	b.middlewares[name] = m
}
func (b *Plugin) Server() {
	plugin.Serve(b.build())
}
func (b *Plugin) build() *plugin.ServeConfig {

	p := &PluginServer{
		httpHandler: b.httpHandler,
		middlewares: b.middlewares,
	}
	if b.httpHandler == nil {
		p.httpHandler = http.NotFoundHandler()
	}
	return &plugin.ServeConfig{
		HandshakeConfig: shared.HandshakeConfig,
		TLSProvider:     nil,
		Plugins: map[string]plugin.Plugin{
			shared.PluginHandlerName: &GrpcPluginServer{
				iml: p,
			},
		},
		VersionedPlugins: nil,
		GRPCServer:       plugin.DefaultGRPCServer,
		Logger:           logger,
		Test:             nil,
	}
}
func ProcessRequestBy(f MiddlewareHandFunc) SetFunc {
	return func(m *serverMiddlewareHandler) {
		m.requestHandler = f
	}
}
func ProcessResponseBy(f MiddlewareHandFunc) SetFunc {
	return func(m *serverMiddlewareHandler) {
		m.responseHandler = f
	}
}

type PluginServer struct {
	proto2.UnimplementedServiceServer
	httpHandler shared.HttpHandler
	middlewares map[string]*serverMiddlewareHandler
}

func (p *PluginServer) MiddlewaresRequest(ctx context.Context, request *proto2.MiddlewareRequest) (*proto2.MiddlewareResponse, error) {

	m, ok := p.middlewares[request.Name]
	if !ok {
		return new(proto2.MiddlewareResponse), nil
	}

	return doMiddlewareHandler(m.requestHandler, ctx, request)
}

func (p *PluginServer) MiddlewaresResponse(ctx context.Context, request *proto2.MiddlewareRequest) (*proto2.MiddlewareResponse, error) {
	response := new(proto2.MiddlewareResponse)
	m, ok := p.middlewares[request.Name]
	if !ok {
		return response, nil
	}

	return doMiddlewareHandler(m.responseHandler, ctx, request)

}
func doMiddlewareHandler(handFunc MiddlewareHandFunc, ctx context.Context, request *proto2.MiddlewareRequest) (*proto2.MiddlewareResponse, error) {
	response := new(proto2.MiddlewareResponse)
	if handFunc == nil {
		return response, nil
	}
	r := &module.MiddlewareRequest{
		FulPath: request.Request.FulPath,
		Url:     request.Request.Url,
		Method:  request.Request.Method,
		Header:  make(map[string][]string),
	}
	for _, h := range request.Request.Headers {
		r.Header[h.Key] = h.Value
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
func (p *PluginServer) GetMiddlewareInfo(ctx context.Context, r *proto2.MiddlewareInfoRequest) (*proto2.MiddlewareInfoResponse, error) {
	info := &proto2.MiddlewareInfoResponse{
		Name:     r.Name,
		Request:  false,
		Response: false,
	}
	m, ok := p.middlewares[r.Name]
	if !ok {
		return info, nil
	}
	info.Request = m.requestHandler != nil
	info.Response = m.responseHandler != nil
	return info, nil
}

func (p *PluginServer) Request(ctx context.Context, request *proto2.HttpRequest) (*proto2.HttpResponse, error) {
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
	iml *PluginServer
}

func (p *GrpcPluginServer) GRPCServer(broker *plugin.GRPCBroker, server *grpc.Server) error {
	proto2.RegisterServiceServer(server, p.iml)
	return nil
}

func (p *GrpcPluginServer) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, conn *grpc.ClientConn) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}
