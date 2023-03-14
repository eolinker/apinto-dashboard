//go:build release

package frontend_controller

import (
	"context"
	"embed"
	_ "embed"
	"fmt"
	"github.com/eolinker/apinto-dashboard/access"
	gzip "github.com/eolinker/apinto-dashboard/common/gzip-static"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/gin-gonic/gin"
	"io/fs"
	"net/http"
	"strings"
	"time"
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
	r.GET("/login", loginHtml)
	r.GET("/auth", authHtml)
	r.NoRoute(noRoute)

}

func authHtml(ginCtx *gin.Context) {
	ginCtx.Header("Cache-Control", "no-store, no-cache, max-age=0, must-revalidate, proxy-revalidate")
	ginCtx.Data(http.StatusOK, "text/html; charset=utf-8", indexContext)
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
	ctx := context.Background()
	session, _ := ginCtx.Cookie(controller.Session)
	modulePath := ginCtx.Request.RequestURI[1:]

	valid, _ := controller.authService.CheckCertValid(ctx)
	if !valid {
		redirect := "/auth"
		ginCtx.Redirect(http.StatusFound, redirect)
		return
	}

	if session != "" {
		_, err := controller.sessionCache.Get(ctx, controller.sessionCache.Key(session))
		if err == nil {

			//此处代码问题 更深层次的路由，比如上游管理之下的上线管理，ModuleMap只记录了一级路由，会导致匹配不到跳转至首页
			//但实际是要跳转到具体某个页面，前端已经做了这部分工作，所以代码可以先注释

			//count := strings.Count(modulePath, "?")
			//tempPath := modulePath
			//if count > 0 {
			//	tempPath = modulePath[:strings.Index(modulePath, "?")]
			//}

			//if _, ok := access.ModuleMap[tempPath]; !ok {
			//	ginCtx.Redirect(http.StatusFound, "/")
			//	return
			//}

			ginCtx.Header("Cache-Control", "no-store, no-cache, max-age=0, must-revalidate, proxy-revalidate")
			ginCtx.Data(http.StatusOK, "text/html; charset=utf-8", indexContext)

			return
		}
	}

	count := strings.Count(modulePath, "?")
	tempPath := modulePath
	if count > 0 {
		tempPath = modulePath[:strings.Index(modulePath, "?")]
	}

	redirect := "/login"
	if _, ok := access.ModuleMap[tempPath]; ok {
		redirect = fmt.Sprint(redirect, "?callback=", ginCtx.Request.RequestURI)
	}
	ginCtx.SetCookie(controller.Session, "", -1, "", "", false, true)

	ginCtx.Redirect(http.StatusFound, redirect)
	return

}
func loginHtml(ginCtx *gin.Context) {
	ctx := context.Background()

	valid, _ := controller.authService.CheckCertValid(ctx)
	if !valid {
		redirect := "/auth"
		ginCtx.Redirect(http.StatusFound, redirect)
		return
	}

	session, _ := ginCtx.Cookie(controller.Session)
	if session != "" {
		_, err := controller.sessionCache.Get(ctx, controller.sessionCache.Key(session))
		if err == nil {
			callback := ginCtx.Query("callback")
			if len(callback) == 0 {
				callback = "/"
			}
			if _, ok := access.ModuleMap[callback[1:]]; !ok {
				callback = "/"
			}
			ginCtx.Redirect(http.StatusFound, callback)
			return
		}

	}
	ginCtx.Header("Cache-Control", "no-store, no-cache, max-age=0, must-revalidate, proxy-revalidate")
	ginCtx.Data(http.StatusOK, "text/html; charset=utf-8", indexContext)
}
