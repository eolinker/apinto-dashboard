//go:build release

package controller

import (
	"embed"
	_ "embed"
	"fmt"
	"io/fs"
	"net/http"
	"path"
	"strings"
	"time"

	gzip "github.com/eolinker/apinto-dashboard/common/gzip-static"
	"github.com/gin-gonic/gin"
)

var (
	expires      = time.Hour * 24 * 7
	cacheControl = fmt.Sprintf("public, max-age=%d", 3600*24*7)
)

type joinPathFunc func(c *gin.Context) string

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
	r.Group("/ace-builds", addExpires, gzip.Gzip(gzip.DefaultCompression)).StaticFS("/", getFileSystem("dist/ace-builds"))
	r.Group("/frontend", addExpires, gzip.Gzip(gzip.DefaultCompression)).StaticFS("/", getFileSystem("dist"))

	r.GET("/favicon.ico", addExpires, func(ginCtx *gin.Context) {
		ginCtx.Writer.WriteString(iconContext)
	})
	r.GET("/", indexHtml)
	r.NoRoute(noRoute)
}

func EmbedPluginFrontend(r *gin.Engine) {
	fileDir := "../../controller/plugin"

	//获取插件图标
	StaticFS(r.Group("/plugin/icon/:id/:file"), "/", "", gin.Dir(fileDir, false), func(c *gin.Context) string {
		id := c.Param("id")
		file := c.Param("file")
		return fmt.Sprintf("%s/%s", id, file)
	})

	//获取插件描述中要用到的MD文件
	StaticFS(r.Group("/plugin/md/:id/:file"), "/", "", gin.Dir(fileDir, false), func(c *gin.Context) string {
		id := c.Param("id")
		file := c.Param("file")
		return fmt.Sprintf("%s/%s", id, file)
	})
	//插件描述中要用到的资源
	StaticFS(r.Group("/plugin/info/:id/resources"), "/", "/*filepath", gin.Dir(fileDir, false), func(c *gin.Context) string {
		id := c.Param("id")
		filePath := c.Param("filepath")
		return fmt.Sprintf("%s/resources/%s", id, filePath)
	})

	//插件详情页
	StaticFS(r.Group("/plugin/info/:id/:file"), "/", "", gin.Dir(fileDir, false), func(c *gin.Context) string {
		id := c.Param("id")
		file := c.Param("file")
		return fmt.Sprintf("%s/%s", id, file)
	})

	//获取插件列表
	r.GET("/plugin/list/:group", indexHtml)
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

func StaticFS(group *gin.RouterGroup, relativePath, extendedPath string, fs http.FileSystem, pathFunc joinPathFunc) {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static folder")
	}
	handler := createStaticHandler(relativePath, fs, pathFunc)
	urlPattern := relativePath
	if extendedPath != "" {
		urlPattern = path.Join(relativePath, extendedPath)
	}

	// Register GET and HEAD handlers
	group.GET(urlPattern, handler)
	group.HEAD(urlPattern, handler)
}

func createStaticHandler(relativePath string, fs http.FileSystem, pathFunc joinPathFunc) gin.HandlerFunc {
	fileServer := http.StripPrefix(relativePath, http.FileServer(fs))

	return func(c *gin.Context) {
		file := pathFunc(c)
		// Check if file exists and/or if we have permission to access it
		f, err := fs.Open(file)
		if err != nil {
			c.Writer.WriteHeader(http.StatusNotFound)
			//TODO 文件不存在时的handler
			//c.handlers = group.engine.noRoute
			// Reset index
			//c.index = -1
			return
		}
		f.Close()

		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}
