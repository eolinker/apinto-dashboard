package local

import (
	"bytes"
	"fmt"
	apinto_module "github.com/eolinker/apinto-module"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"io"
	"net/http"
	"strings"
)

const (
	MiddlewarePre  = "pre"
	MiddlewarePost = "post"
)

func (p *ProxyAPi) CreateMiddleware(name, path, life string, rule [][]string) apinto_module.MiddlewareHandler {
	url := fmt.Sprintf("%s/%s", p.server, strings.TrimPrefix(path, "/"))
	switch life {
	case MiddlewarePost:
		return apinto_module.MiddlewareHandler{
			Name: name,
			Rule: apinto_module.CreateMiddlewareRules(rule),
			Handler: func(ginCtx *gin.Context) {
				ginCtx.Next()
				doMiddleware(ginCtx, url, p.module)
			},
		}
	case MiddlewarePre:

	}
	return apinto_module.MiddlewareHandler{
		Name: name,
		Rule: apinto_module.CreateMiddlewareRules(rule),
		Handler: func(ginCtx *gin.Context) {

			doMiddleware(ginCtx, url, p.module)
		},
	}
}
func doMiddleware(ginCtx *gin.Context, url string, module string) {
	request := apinto_module.CreateMiddlewareRequest(ginCtx, module)
	marshal, _ := json.Marshal(request)
	buf := bytes.NewReader(marshal)
	response, err := http.Post(url, "application/json", buf)
	if err != nil {

		ginCtx.AbortWithError(500, err)
		return
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		ginCtx.AbortWithError(500, err)
		return
	}
	response.Body.Close()

	middlewareResponse, err := apinto_module.DecodeMiddlewareResponse(body)
	if err != nil {
		return
	}
	if len(middlewareResponse.Header) > 0 {
		for k, v := range middlewareResponse.Header {
			ginCtx.Header(k, v)
		}
	}
	if len(middlewareResponse.Keys) > 0 {
		for k, v := range middlewareResponse.Keys {
			ginCtx.Set(k, v)
		}
	}
	if middlewareResponse.Abort {
		if middlewareResponse.Body != nil {
			ginCtx.Data(middlewareResponse.StatusCode, middlewareResponse.ContentType, middlewareResponse.Body)
		}
		ginCtx.Abort()
	}

}
