package frontend

import (
	"embed"
	"fmt"
	"github.com/eolinker/apinto-dashboard/common/gzip-static"
	"github.com/gin-gonic/gin"
	"io/fs"
	"net/http"
	"time"
)

var (
	//go:embed dist/index.html
	indexContent []byte
	//go:embed dist
	dist embed.FS
	//go:embed dist/favicon.ico
	iconContent []byte
)

var (
	expires      = time.Hour * 24 * 7
	cacheControl = fmt.Sprintf("public, max-age=%d", 3600*24*7)
)
var (
	GzipHandler = gzip.Gzip(gzip.DefaultCompression)
)

func AddExpires(ginCtx *gin.Context) {
	ginCtx.Header("Expires", time.Now().Add(expires).UTC().Format(http.TimeFormat))
	ginCtx.Header("Cache-Control", cacheControl)
}
func GetFileSystem(dir string) http.FileSystem {
	fDir, err := fs.Sub(dist, dir)
	if err != nil {
		panic(err)
	}

	return http.FS(fDir)

}
func Favicon(ginCtx *gin.Context) {
	ginCtx.Data(http.StatusOK, "image/x-icon", iconContent)
}

func IndexHtml(ginCtx *gin.Context) {
	ginCtx.Header("Cache-Control", "no-store, no-cache, max-age=0, must-revalidate, proxy-revalidate")
	ginCtx.Data(http.StatusOK, "text/html; charset=utf-8", indexContent)
	return
}
