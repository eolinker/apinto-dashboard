package controller

import (
	"embed"
	apinto_module "github.com/eolinker/apinto-module"
	"github.com/gin-gonic/gin"
	"io/fs"
	"net/http"
	"path"
)

//go:embed dist/index.html
var indexContext []byte

//go:embed dist
var dist embed.FS

//go:embed dist/favicon.ico
var iconContext []byte

func IndexHtml(ginCtx gin.Context) {
	ginCtx.Header("Cache-Control", "no-store, no-cache, max-age=0, must-revalidate, proxy-revalidate")
	ginCtx.Data(http.StatusOK, "text/html; charset=utf-8", indexContext)
	return
}

func favicon() apinto_module.RouterInfo {
	return apinto_module.RouterInfo{
		Method:  http.MethodGet,
		Path:    "/favicon.ico",
		Handler: "favicon.ico",
		Labels:  apinto_module.RouterLabelAssets,
		HandlerFunc: []apinto_module.HandlerFunc{func(ginCtx *gin.Context) {
			ginCtx.Data(http.StatusOK, "image/x-icon", iconContext)
		}},
	}
}

func getFileSystem(dir string) http.FileSystem {
	fDir, err := fs.Sub(dist, dir)
	if err != nil {
		panic(err)
	}

	return http.FS(fDir)

}

func staticFile(prefix string, dir string) apinto_module.RoutersInfo {
	fileSystem := getFileSystem(dir)
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
			HandlerFunc: []apinto_module.HandlerFunc{handler},
		},
		{
			Method:      http.MethodHead,
			Path:        path.Join(prefix, "/*filepath"),
			Handler:     prefix,
			Labels:      apinto_module.RouterLabelAssets,
			HandlerFunc: []apinto_module.HandlerFunc{handler},
		},
	}
}
