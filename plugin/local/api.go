package local

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/eolinker/apinto-dashboard/controller"
	apinto_module "github.com/eolinker/apinto-module"
	"github.com/eolinker/eosc/log"
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
	if config.Path != "" {
		to = config.Path
	}
	return p.createApi(name, method, fmt.Sprintf("/api/module//%s/%s", p.module, strings.TrimPrefix(path, "/")), to, mergeLabel(apinto_module.RouterLabelApi, config.Label))
}
func (p *ProxyAPi) CreateOpenApi(name, method, path string, config PathConfig) apinto_module.RouterInfo {
	to := path
	if config.Path != "" {
		to = config.Path
	}
	return p.createApi(name, method, fmt.Sprintf("/api2/module/%s/%s", p.module, strings.TrimPrefix(path, "/")), to, mergeLabel(apinto_module.RouterLabelOpenApi, config.Label))
}
func (p *ProxyAPi) createApi(name, method, from, to string, labels []string) apinto_module.RouterInfo {

	return apinto_module.RouterInfo{
		Method:  method,
		Path:    from,
		Handler: fmt.Sprintf("%s.%s", p.module, name),
		Labels:  labels,
		HandlerFunc: []apinto_module.HandlerFunc{func(ginCtx *gin.Context) {

			targetPath := to
			for _, param := range ginCtx.Params {
				targetPath = strings.Replace(targetPath, fmt.Sprint(":", param.Key), param.Value, -1)
			}
			query := ginCtx.Request.URL.Query()
			for k, v := range p.query {
				query.Set(k, v)
			}

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

			headerName, value := apinto_module.ReadKeys(ginCtx)
			if headerName != "" {
				request.Header.Set(headerName, value)
			}
			response, err := http.DefaultClient.Do(request)
			if err != nil {
				return
			}
			responseData, _ := io.ReadAll(response.Body)
			response.Body.Close()

			responseHeader := response.Header
			contentType := response.Header.Get("content-type")
			for k, vs := range responseHeader {
				switch k {
				case "Apinto-Event":
					for _, v := range vs {
						if len(v) > 0 {
							eventObjs := make(map[string]any)
							if err := json.Unmarshal([]byte(v), &eventObjs); err == nil {
								for k, v := range eventObjs {
									apinto_module.DoEvent(k, v)
								}
							} else {
								log.Warnf("invalid event for :%s.%s on %s %s", p.module, name, method, to)
							}
						}
					}
				default:
					for _, v := range vs {
						ginCtx.Writer.Header().Add(k, v)
					}

				}
			}
			ginCtx.Data(response.StatusCode, contentType, responseData)
		}},
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
