package local

import (
	"bytes"
	"fmt"
	"github.com/eolinker/apinto-dashboard/controller"
	apinto_module "github.com/eolinker/apinto-dashboard/module"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
)

type ProxyAPi struct {
	server  string
	module  string
	headers map[string]string
	query   map[string]string
}

func NewProxyAPi(server string, module string, config *Config) *ProxyAPi {
	if !strings.HasPrefix(server, "http://") && !strings.HasPrefix(server, "https://") {
		server = fmt.Sprintf("http://%s", server)
	}
	server = strings.TrimSuffix(server, "/")
	return &ProxyAPi{server: server, module: module, headers: config.Header, query: config.Query}
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
func (p *ProxyAPi) proxyApiHandler(name, method, path string) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		targetPath := path
		for _, param := range ginCtx.Params {
			targetPath = strings.Replace(targetPath, fmt.Sprint(":", param.Key), param.Value, -1)
			targetPath = strings.Replace(targetPath, fmt.Sprint("*", param.Key), param.Value, -1)
		}
		query := ginCtx.Request.URL.Query()
		for k, v := range p.query {
			query.Set(k, v)
		}
		targetPath = strings.TrimPrefix(targetPath, "/")
		url := fmt.Sprintf("%s/%s?%s", p.server, targetPath, query.Encode())
		data, err := ginCtx.GetRawData()
		if err != nil {
			return
		}
		request, err := http.NewRequest(method, url, bytes.NewReader(data))
		if err != nil {
			controller.ErrorJson(ginCtx, 200, err.Error())
			return
		}

		request.Header = ginCtx.Request.Header
		for k, v := range p.headers {
			request.Header.Set(k, v)
		}
		apinto_module.CopyKeysToHeader(ginCtx.Keys, request)

		response, err := http.DefaultClient.Do(request)
		if err != nil {
			return
		}
		responseData, _ := io.ReadAll(response.Body)
		response.Body.Close()

		apinto_module.ReSetHeaderFromProxyResponse(response.Header, ginCtx)
		contentType := response.Header.Get("content-type")
		ginCtx.Data(response.StatusCode, contentType, responseData)
	}
}
func (p *ProxyAPi) createApi(name, method, from, to string, labels []string) apinto_module.RouterInfo {

	return apinto_module.RouterInfo{
		Method:      method,
		Path:        from,
		Handler:     fmt.Sprintf("%s.%s", p.module, name),
		Labels:      labels,
		HandlerFunc: []apinto_module.HandlerFunc{p.proxyApiHandler(name, method, to)},
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
