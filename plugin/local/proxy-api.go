package local

import (
	"context"
	"fmt"
	apinto_module "github.com/eolinker/apinto-dashboard/module"
	client "github.com/eolinker/apinto-dashboard/plugin/go-plugin/plugin-client"
	"github.com/eolinker/apinto-dashboard/plugin/go-plugin/proto"
	"github.com/eolinker/apinto-dashboard/plugin/go-plugin/shared"
	"github.com/eolinker/eosc/log"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-plugin"
	"strings"
	"sync"
)

type ProxyAPi struct {
	client   *plugin.Client
	cmd      string
	module   string
	params   map[string]string
	killOnce sync.Once

	wg sync.WaitGroup

	initOnce    sync.Once
	handler     client.ClientHandler
	clientError error
}

func (p *ProxyAPi) kill() {
	p.killOnce.Do(func() {
		p.wg.Done()
		p.wg.Wait()
		if p.client != nil {
			p.client.Kill()
			p.client = nil
		}
	})
}
func (p *ProxyAPi) initClient() {
	p.initOnce.Do(func() {
		params := make([]string, 0, len(p.params))
		for k, v := range p.params {
			params = append(params, fmt.Sprintf("%s=%s", k, v))
		}
		c := client.CreateClient(p.module, cmdPath(p.cmd), params...)

		rpcClient, err := c.Client()
		if err != nil {
			p.clientError = err
			return
		}
		// Request the plugin
		raw, err := rpcClient.Dispense(shared.PluginHandlerName)
		if err != nil {
			p.clientError = err
			return
		}

		// We should have a Greeter now! This feels like a normal interface
		// implementation but is in fact over an RPC connection.
		handler := raw.(client.ClientHandler)
		p.handler = handler
		p.client = c
	})
}
func NewProxyAPi(cmd string, module string, config *Config) *ProxyAPi {

	p := &ProxyAPi{cmd: cmd, module: module, params: config.Initialize, wg: sync.WaitGroup{}}
	p.wg.Add(1)
	return p
}

func (p *ProxyAPi) CreateApi(name, method, path string, config PathConfig) apinto_module.RouterInfo {
	to := path
	routerPath := config.Path
	if routerPath == "" {
		routerPath = fmt.Sprintf("/api/module/%s/%s", p.module, strings.TrimPrefix(path, "/"))
	}
	label := config.Label
	if len(label) == 0 {
		label = apinto_module.RouterLabelApi
	}
	return p.createApi(name, method, routerPath, to, label)
}
func (p *ProxyAPi) CreateOpenApi(name, method, path string, config PathConfig) apinto_module.RouterInfo {
	to := path
	routerPath := config.Path
	if routerPath == "" {
		routerPath = fmt.Sprintf("/ap2/module/%s/%s", p.module, strings.TrimPrefix(path, "/"))
	}
	label := config.Label
	if len(label) == 0 {
		label = apinto_module.RouterLabelApi
	}
	return p.createApi(name, method, routerPath, to, label)
}
func (p *ProxyAPi) proxyApiHandler(from, to string) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {

		p.initClient()
		if p.clientError != nil {
			log.Error(ginCtx.AbortWithError(500, p.clientError))
			return
		}

		p.wg.Add(1)
		defer p.wg.Done()

		targetPath := to
		sourcePath := from
		for _, param := range ginCtx.Params {
			targetPath = strings.Replace(targetPath, fmt.Sprint(":", param.Key), param.Value, -1)
			targetPath = strings.Replace(targetPath, fmt.Sprint("*", param.Key), param.Value, -1)

			sourcePath = strings.Replace(targetPath, fmt.Sprint(":", param.Key), param.Value, -1)
			sourcePath = strings.Replace(targetPath, fmt.Sprint("*", param.Key), param.Value, -1)
		}

		targetPath = strings.TrimPrefix(targetPath, "/")
		sourcePath = strings.TrimPrefix(sourcePath, "/")
		p.handler.ServerGin(ginCtx, func(ctx context.Context, request *proto.HttpRequest) {
			request.Url = strings.Replace(ginCtx.Request.URL.String(), sourcePath, targetPath, -1)
		}, nil)

	}
}
func (p *ProxyAPi) createApi(name, method, from, to string, labels []string) apinto_module.RouterInfo {

	return apinto_module.RouterInfo{
		Method:      method,
		Path:        from,
		Handler:     fmt.Sprintf("%s.%s", p.module, name),
		Labels:      labels,
		HandlerFunc: []apinto_module.HandlerFunc{p.proxyApiHandler(from, to)},
	}
}

func mergeLabel[T comparable](args ...[]T) []T {

	m := make(map[T]struct{})
	for _, ls := range args {
		for _, v := range ls {
			m[v] = struct{}{}
		}
	}
	rs := make([]T, 0, len(m))
	for k := range m {
		rs = append(rs, k)
	}
	return rs
}
