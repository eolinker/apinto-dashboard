package controller

import (
	"github.com/eolinker/apinto-dashboard/frontend"
	apinto_module "github.com/eolinker/apinto-module"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
)

func indexRouter() apinto_module.RouterInfo {
	return apinto_module.RouterInfo{
		Method:      http.MethodGet,
		Path:        "/",
		Handler:     "index",
		Labels:      apinto_module.RouterLabelModule,
		HandlerFunc: []apinto_module.HandlerFunc{frontend.IndexHtml},
	}
}
func favicon() apinto_module.RouterInfo {
	return apinto_module.RouterInfo{
		Method:      http.MethodGet,
		Path:        "/favicon.ico",
		Handler:     "favicon.ico",
		Labels:      apinto_module.RouterLabelAssets,
		HandlerFunc: []apinto_module.HandlerFunc{frontend.Favicon},
	}
}

func staticFile(prefix string, dir string) apinto_module.RoutersInfo {
	fileSystem := frontend.GetFileSystem(dir)
	handler := func(ginCtx *gin.Context) {
		filePath := ginCtx.Param("filepath")
		ginCtx.FileFromFS(filePath, fileSystem)

	}
	return apinto_module.RoutersInfo{
		{
			Method:      http.MethodGet,
			Path:        path.Join(prefix, "/*filepath"),
			Handler:     prefix,
			Labels:      apinto_module.RouterLabelAssets,
			HandlerFunc: []apinto_module.HandlerFunc{frontend.AddExpires, frontend.GzipHandler, handler},
		},
		{
			Method:      http.MethodHead,
			Path:        path.Join(prefix, "/*filepath"),
			Handler:     prefix,
			Labels:      apinto_module.RouterLabelAssets,
			HandlerFunc: []apinto_module.HandlerFunc{frontend.AddExpires, frontend.GzipHandler, handler},
		},
	}
}
