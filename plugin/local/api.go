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
func (p *ProxyAPi) CreateApi(name, method, path string, appendLabel []string) apinto_module.RouterInfo {
	return p.createApi(name, method, fmt.Sprintf("/api/module/%s", p.module), path, mergeLabel(apinto_module.RouterLabelApi, appendLabel))
}
func (p *ProxyAPi) CreateOpenApi(name, method, path string, appendLabel []string) apinto_module.RouterInfo {
	return p.createApi(name, method, fmt.Sprintf("/api2/module/%s", p.module), path, mergeLabel(apinto_module.RouterLabelOpenApi, appendLabel))
}
func (p *ProxyAPi) createApi(name, method, prefix, path string, labels []string) apinto_module.RouterInfo {
	prefix = strings.Trim(prefix, "/")
	path = strings.TrimPrefix(path, "/")
	routerPath := fmt.Sprintf("/%s/%s", prefix, strings.TrimPrefix(path, "/"))
	pathPrefix := fmt.Sprintf("/%s/", prefix)
	return apinto_module.RouterInfo{
		Method:  method,
		Path:    routerPath,
		Handler: fmt.Sprintf("%s.%s", p.module, name),
		Labels:  labels,
		HandlerFunc: []apinto_module.HandlerFunc{func(ginCtx *gin.Context) {
			requestPath := ginCtx.Request.RequestURI
			targetPath := strings.TrimPrefix(requestPath, pathPrefix)
			url := fmt.Sprintf("%s/%s", p.server, targetPath)
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
			if ginCtx.Keys != nil {
				apintoKeysData, err := json.Marshal(ginCtx.Keys)
				if err == nil {
					request.Header.Set("apinto-runtime-keys", string(apintoKeysData))
				}
			}

			response, err := http.DefaultClient.Do(request)
			if err != nil {
				return
			}
			responseData, _ := io.ReadAll(response.Body)
			response.Body.Close()

			event := response.Header.Get("apinto-event")

			contentType := response.Header.Get("content-type")

			if len(event) > 0 {
				eventObjs := make(map[string]any)
				if err := json.Unmarshal([]byte(event), &eventObjs); err == nil {
					for k, v := range eventObjs {
						apinto_module.DoEvent(k, v)
					}
				} else {
					log.Warnf("invalid event for :%s.%s on %s %s", p.module, name, method, path)
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
