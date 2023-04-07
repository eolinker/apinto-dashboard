//go:build release

package controller

import (
	"embed"
	_ "embed"
	"errors"
	"fmt"
	"github.com/eolinker/eosc/log"
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

	filePathErr = errors.New("filePath is illegal. ")
)

type joinPathFunc func(c *gin.Context) (string, error)

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
	fileDir := "./plugin"
	pluginFs := gin.Dir(fileDir, false)

	//获取插件图标
	StaticFS(r.Group("/plugin/icon"), "/plugin/icon", pluginFs, func(c *gin.Context) (string, error) {
		filePath := strings.Trim(c.Param("filepath"), "/")
		list := strings.Split(filePath, "/")
		pathLen := len(list)
		//若文件路径层级多于 id/file 或 filepath 为空
		if pathLen > 2 || filePath == "" {
			return "", filePathErr
		}
		//防止读取resources目录
		if pathLen == 2 && list[1] == "resources" {
			return "", filePathErr
		}
		if pathLen == 1 {
			//TODO 从插件配置文件获取图标名
			fileName := ""
			return fmt.Sprintf("%s/%s", list[0], fileName), nil
		}
		return filePath, nil
	})

	//获取插件描述中要用到的MD文件
	StaticFS(r.Group("/plugin/md"), "/plugin/md", pluginFs, func(c *gin.Context) (string, error) {
		filePath := strings.Trim(c.Param("filepath"), "/")
		list := strings.Split(filePath, "/")
		pathLen := len(list)
		//若文件路径层级多于 id/file 或 filepath 为空
		if pathLen > 2 || filePath == "" {
			return "", filePathErr
		}
		//防止读取resources目录
		if pathLen == 2 && list[1] == "resources" {
			return "", filePathErr
		}
		if pathLen == 1 {
			return fmt.Sprintf("%s/%s", list[0], "README.md"), nil
		}
		return filePath, nil
	})

	//插件详情页  插件描述中要用到的资源
	StaticFS(r.Group("/plugin/info"), "/plugin/info", pluginFs, func(c *gin.Context) (string, error) {
		filePath := strings.Trim(c.Param("filepath"), "/")

		list := strings.Split(filePath, "/")
		pathLen := len(list)
		//若filepath为空
		if filePath == "" {
			return "", filePathErr
		}
		//防止读取resources目录
		if pathLen == 2 && list[1] == "resources" {
			return "", filePathErr
		}
		//若文件路径层级为 id/resources/*filepath
		if pathLen > 2 && list[1] == "resources" {
			return filePath, nil
		}
		if pathLen == 1 {
			return fmt.Sprintf("%s/%s", list[0], "README.md"), nil
		}

		return filePath, nil
	})

	//获取插件列表
	r.GET("/plugin/list/:group", indexHtml)
}
func favicon(ginCtx *gin.Context) {
	ginCtx.Writer.WriteString(iconContext)
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

func StaticFS(group *gin.RouterGroup, relativePath string, fs http.FileSystem, pathFunc joinPathFunc) {
	handler := createStaticHandler(relativePath, fs, pathFunc)

	// Register GET and HEAD handlers
	group.GET("/*filepath", handler)
	group.HEAD("/*filepath", handler)
}

func createStaticHandler(relativePath string, fs http.FileSystem, pathFunc joinPathFunc) gin.HandlerFunc {
	fileServer := http.StripPrefix(relativePath, http.FileServer(fs))

	return func(c *gin.Context) {
		file, err := pathFunc(c)
		if err != nil {
			log.Info(fmt.Errorf("%s:%w", c.Request.URL.Path, err))
			c.Data(http.StatusNotFound, "application/text", []byte("404 page not found"))
			return
		}
		// Check if file exists and/or if we have permission to access it
		f, err := fs.Open(file)
		if err != nil {
			//c.Writer.WriteHeader(http.StatusNotFound)
			//TODO 文件不存在时
			//c.handlers = group.engine.noRoute
			// Reset index
			//c.index = -1
			c.Data(http.StatusNotFound, "application/text", []byte("404 page not found"))
			return
		}
		f.Close()

		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}
