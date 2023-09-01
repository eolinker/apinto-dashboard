package plugin

import (
	"context"
	"github.com/gin-gonic/gin"
	"os"
)

type ginKeysKey struct{}

var (
	GinKeys = &ginKeysKey{}
)

func Engine() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.LoggerWithWriter(os.Stderr))
	engine.Use(KeyForGin)
	engine.ContextWithFallback = true
	return engine
}
func widthKey(ctx context.Context, keys map[string]interface{}) context.Context {
	ctx = context.WithValue(ctx, GinKeys, keys)
	return ctx
}
func KeyForGin(ginCtx *gin.Context) {
	value := ginCtx.Value(GinKeys)
	if value != nil {
		keys := value.(map[string]interface{})
		if keys == nil {
			return
		}
		for k, v := range keys {
			ginCtx.Set(k, v)
		}
	}
}
