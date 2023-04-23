package local

import (
	"bytes"
	"fmt"
	apinto_module "github.com/eolinker/apinto-module"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
)

func (p *ProxyAPi) CreateHome(path string) apinto_module.RouterInfo {
	baseHtml := []byte(fmt.Sprintf(fmt.Sprintf("<base href=\"/module/%s\">", p.module)))
	routerRoot := fmt.Sprintf("/module/%s/", p.module)
	routerPath := routerRoot

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if strings.HasSuffix(path, "/") {
		routerPath = fmt.Sprintf("/module/%s/*path", p.module)
	}
	return apinto_module.RouterInfo{
		Method:  http.MethodGet,
		Path:    routerPath,
		Handler: fmt.Sprintf("%s.home", p.module),
		Labels:  apinto_module.RouterLabelModule,
		HandlerFunc: []apinto_module.HandlerFunc{func(ginCtx *gin.Context) {
			query := ginCtx.Request.URL.Query()
			for k, v := range p.query {
				query.Set(k, v)
			}
			header := ginCtx.Request.Header
			for k, v := range p.headers {
				header.Set(k, v)
			}

			targetPath := strings.Replace(ginCtx.Request.URL.Path, routerRoot, path, 1)

			target := fmt.Sprintf("%s%s?%s", p.server, targetPath, query.Encode())

			request, err := http.NewRequest(http.MethodGet, target, nil)
			if err != nil {
				ginCtx.Error(err)
				return
			}
			request.Header = header
			response, err := http.DefaultClient.Do(request)
			if err != nil {
				ginCtx.Error(err)
				return
			}
			body, err := io.ReadAll(response.Body)
			response.Body.Close()
			if err != nil {
				ginCtx.Error(err)
				return
			}
			statusCode := response.StatusCode
			contentType := response.Header.Get("content-type")

			if bytes.Index(body, []byte(`<base href="/">`)) > 0 {
				body = bytes.Replace(body, []byte(`<base href="/">`), baseHtml, 1)
			} else {
				body = bytes.Replace(body, []byte(`</head>`), bytes.Join([][]byte{baseHtml, []byte(`</head>`)}, []byte("\n")), 1)
			}

			ginCtx.Data(statusCode, contentType, body)
		}},
	}
}
func (p *ProxyAPi) CreateHtml(dir string, appendLabel []string) apinto_module.RouterInfo {
	routerRoot := fmt.Sprintf("/module/%s/%s", p.module, strings.TrimPrefix(dir, "/"))
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
			path := strings.Replace(ginCtx.Request.RequestURI, routerRoot, targetPath, 1)
			url := fmt.Sprintf("%s%s", p.server, path)
			response, err := http.Get(url)
			if err != nil {
				ginCtx.Error(err)
				return
			}
			body, err := io.ReadAll(response.Body)
			response.Body.Close()
			if err != nil {
				ginCtx.Error(err)
				return
			}
			if response.Header.Get("context-type")
			for k, vs := range response.Header {
				if len(vs) > 0 {
					ginCtx.Header(k, vs[0])
				} else {
					ginCtx.Header(k, "")
				}
			}

			ginCtx.Writer.WriteHeader(response.StatusCode)
			ginCtx.Writer.Write(body)
		}},
	}
}
