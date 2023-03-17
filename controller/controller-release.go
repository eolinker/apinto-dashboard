//go:build release

package controller

import (
	"embed"
	_ "embed"
	"fmt"
	"io/fs"
	"net/http"
	"strings"
	"time"

	gzip "github.com/eolinker/apinto-dashboard/common/gzip-static"
	"github.com/gin-gonic/gin"
)

var (
	expires      = time.Hour * 24 * 7
	cacheControl = fmt.Sprintf("public, max-age=%d", 3600*24*7)
)

//go:embed dist/index.html
var indexContext []byte

//go:embed dist
var dist embed.FS

//go:embed dist/favicon.ico
var iconContext string

func getFileSystem(dir string) http.FileSystem {
	fs, err := fs.Sub(dist, dir)
	if err != nil {
		panic(err)
	}

	return http.FS(fs)
}
func EmbedFrontend(r *gin.Engine) {

	r.Group("/assets", addExpires, gzip.Gzip(gzip.DefaultCompression)).StaticFS("/", getFileSystem("dist/assets"))
	r.Group("/frontend", addExpires, gzip.Gzip(gzip.DefaultCompression)).StaticFS("/", getFileSystem("dist"))

	r.GET("/favicon.ico", addExpires, func(ginCtx *gin.Context) {
		ginCtx.Writer.WriteString(iconContext)
	})
	r.GET("/", indexHtml)
	r.NoRoute(noRoute)

}

func noRoute(ginCtx *gin.Context) {
	if strings.HasPrefix(ginCtx.Request.URL.Path, "/api/") {
		ginCtx.Data(http.StatusNotFound, "application/json", []byte("404 page not found"))
	} else if strings.Count(ginCtx.Request.URL.Path, ".") > 0 {
		ginCtx.Data(http.StatusNotFound, "application/text", []byte("404 page not found"))
	} else {
		indexHtml(ginCtx)
	}
}

func addExpires(ginCtx *gin.Context) {
	ginCtx.Header("Expires", time.Now().Add(expires).UTC().Format(http.TimeFormat))
	ginCtx.Header("Cache-Control", cacheControl)
}
func indexHtml(ginCtx *gin.Context) {
	ginCtx.Header("Cache-Control", "no-store, no-cache, max-age=0, must-revalidate, proxy-revalidate")
	ginCtx.Data(http.StatusOK, "text/html; charset=utf-8", indexContext)
	return

}
