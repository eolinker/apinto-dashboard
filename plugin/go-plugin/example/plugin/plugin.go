package main

import (
	"context"
	module "github.com/eolinker/apinto-dashboard/module"
	"github.com/eolinker/apinto-dashboard/plugin/go-plugin/plugin"
	"github.com/gin-gonic/gin"
)

func main() {

	engine := plugin.Engine()
	engine.Any("/*all", func(ginCtx *gin.Context) {

		ginCtx.JSON(200, map[string]interface{}{
			"path":    ginCtx.FullPath(),
			"url":     ginCtx.Request.URL.String(),
			"method":  ginCtx.Request.Method,
			"keys":    ginCtx.Keys,
			"headers": ginCtx.Request.Header,
		})
	})
	ps := plugin.NewPlugin(engine.Handler())
	ps.AddMiddleware("test",
		plugin.ProcessRequestBy(func(ctx context.Context, request *module.MiddlewareRequest, writer module.MiddlewareResponseWriter) {
			writer.AddHeader("test", "test")
			writer.Set("test", "test")

		}),
		plugin.ProcessResponseBy(func(ctx context.Context, request *module.MiddlewareRequest, writer module.MiddlewareResponseWriter) {

		}))
	ps.Server()

}
