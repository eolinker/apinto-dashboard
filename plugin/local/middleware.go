package local

import (
	"bytes"
	"fmt"
	apinto_module "github.com/eolinker/apinto-dashboard/module"
	"github.com/eolinker/eosc/log"
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
				doMiddleware(ginCtx, url)
			},
		}
	case MiddlewarePre:

	}
	return apinto_module.MiddlewareHandler{
		Name: name,
		Rule: apinto_module.CreateMiddlewareRules(rule),
		Handler: func(ginCtx *gin.Context) {
			log.Debug("middleware proxy:", ginCtx.Request.RequestURI)
			doMiddleware(ginCtx, url)

		},
	}
}
func doMiddleware(ginCtx *gin.Context, url string) {
	request := apinto_module.CreateMiddlewareRequest(ginCtx)
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
		for k, vs := range middlewareResponse.Header {
			if len(vs) == 1 {
				ginCtx.Header(k, vs[0])
			} else {
				ginCtx.Writer.Header().Del(k)
				for _, v := range vs {
					ginCtx.Writer.Header().Add(k, v)
				}
			}

		}
	}
	if len(middlewareResponse.Keys) > 0 {
		for k, v := range middlewareResponse.Keys {
			ginCtx.Set(k, v)
		}
	}
	if middlewareResponse.Abort {
		ginCtx.Abort()
		switch middlewareResponse.Action {
		case "redirect":
			ginCtx.Redirect(middlewareResponse.StatusCode, string(middlewareResponse.Body))
		default:

			ginCtx.Data(middlewareResponse.StatusCode, middlewareResponse.ContentType, middlewareResponse.Body)

		}

	}

}
