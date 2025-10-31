package frontend

import (
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/eolinker/apinto-dashboard/common/gzip-static"
	"github.com/eolinker/apinto-dashboard/custom"
	apinto_module "github.com/eolinker/apinto-dashboard/module"
	"github.com/eolinker/apinto-dashboard/pm3"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
)

var (
	expires      = time.Hour * 24 * 7
	cacheControl = fmt.Sprintf("public, max-age=%d", 3600*24*7)

	logos = make(map[string]*fileInfo)
)

type fileInfo struct {
	contentType string
	file        []byte
}

func (f *fileInfo) H(ginCtx *gin.Context) {
	ginCtx.Data(http.StatusOK, f.contentType, f.file)
}

func newFileInfo(file []byte) *fileInfo {
	contentType := mimetype.Detect(file).String()
	return &fileInfo{file: file, contentType: contentType}
}

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

func IndexHtml(ginCtx *gin.Context) {

	ginCtx.Header("Cache-Control", "no-store, no-cache, max-age=0, must-revalidate, proxy-revalidate")
	ginCtx.Data(http.StatusOK, "text/html; charset=utf-8", indexContent)
	return
}
func Frontends() []pm3.FrontendAsset {
	//{
	//
	//Path:        "/favicon.ico",
	//	HandlerFunc: favicon,
	//}, {
	//Path:        "/logo.svg",
	//	HandlerFunc: logo,
	//}, {
	//Path: "/logo-auth.svg",
	//},

	fal := make([]pm3.FrontendAsset, 0, 6)
	for p, v := range logos {
		fal = append(fal, pm3.FrontendAsset{
			Path:        p,
			HandlerFunc: v.H,
		})
	}
	return append(fal,
		staticFile("/assets/", "dist/assets"),
		staticFile("/ace-builds/", "dist/ace-builds"),
		staticFile("/frontend/", "dist"),
	)
}
func staticFile(prefix string, dir string) pm3.FrontendAsset {
	fileSystem := GetFileSystem(dir)

	return apinto_module.StaticFileDir(prefix, fileSystem)

}

func init() {
	if d, h := custom.AssetsReplace("logo"); h {
		logos["logo.png"] = newFileInfo(d)
	} else {
		logos["logo.png"] = newFileInfo(logoContent)
	}
	if d, h := custom.AssetsReplace("favicon"); h {
		logos["favicon.ico"] = newFileInfo(d)
	} else {
		logos["favicon.ico"] = newFileInfo(iconContent)
	}

	if d, h := custom.AssetsReplace("auth"); h {
		logos["logo-auth.png"] = newFileInfo(d)
	}

}
