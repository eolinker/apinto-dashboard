package local

import (
	"bytes"
	"context"
	"fmt"
	apinto_module "github.com/eolinker/apinto-dashboard/module"
	"github.com/eolinker/apinto-dashboard/plugin/go-plugin/proto"
	"github.com/eolinker/eosc/log"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func (p *ProxyAPi) CreateHome(path string) []apinto_module.RouterInfo {
	baseHtml := []byte(fmt.Sprintf(fmt.Sprintf("<base href=\"/agent/%s/\">", p.module)))
	routerRoot := fmt.Sprintf("/agent/%s", p.module)

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	handler := func(ginCtx *gin.Context) {
		p.initClient()
		if p.clientError != nil {
			log.Error(ginCtx.AbortWithError(500, p.clientError))
			return
		}
		p.handler.ServerGin(ginCtx, func(ctx context.Context, request *proto.HttpRequest) {
			request.Url = strings.Replace(ginCtx.Request.URL.String(), routerRoot, path, 1)
		}, func(ctx context.Context, response *proto.HttpResponse) {
			if bytes.Index(response.Body, []byte(`<base href="/">`)) > 0 {
				response.Body = bytes.Replace(response.Body, []byte(`<base href="/">`), baseHtml, 1)
			} else {
				response.Body = bytes.Replace(response.Body, []byte(`<head>`), bytes.Join([][]byte{[]byte(`<head>`), baseHtml}, []byte("\n")), 1)
			}
		})

	}
	return []apinto_module.RouterInfo{{
		Method:      http.MethodGet,
		Path:        fmt.Sprintf("/agent/%s", p.module),
		Handler:     fmt.Sprintf("%s.home", p.module),
		Labels:      apinto_module.RouterLabelModule,
		HandlerFunc: []apinto_module.HandlerFunc{handler},
	},
		{
			Method:      http.MethodGet,
			Path:        fmt.Sprintf("/agent/%s/:sub/*path", p.module),
			Handler:     fmt.Sprintf("%s.home", p.module),
			Labels:      apinto_module.RouterLabelModule,
			HandlerFunc: []apinto_module.HandlerFunc{handler},
		}}
}
func (p *ProxyAPi) CreateHtml(dir string, appendLabel []string) apinto_module.RouterInfo {
	routerRoot := fmt.Sprintf("/agent/%s/%s", p.module, strings.TrimPrefix(dir, "/"))
	routerPath := routerRoot
	if strings.HasSuffix(dir, "/") {
		routerPath = fmt.Sprintf("%s*filePath", routerPath)
	}
	targetPath := fmt.Sprintf("/%s", strings.TrimPrefix(dir, "/"))
	return apinto_module.RouterInfo{
		Method:  http.MethodGet,
		Path:    routerPath,
		Handler: fmt.Sprintf("%s.%s", p.module, strings.Join(strings.Split("dir", "/"), ".")),
		Labels:  mergeLabel(apinto_module.RouterLabelAssets, appendLabel),
		HandlerFunc: []apinto_module.HandlerFunc{func(ginCtx *gin.Context) {
			p.initClient()
			if p.clientError != nil {
				log.Error(ginCtx.AbortWithError(500, p.clientError))
				return
			}
			p.wg.Add(1)
			defer p.wg.Done()
			p.handler.ServerGin(ginCtx, func(ctx context.Context, request *proto.HttpRequest) {
				request.Url = strings.Replace(ginCtx.Request.RequestURI, routerRoot, targetPath, 1)
			}, nil)

		}},
	}
}
