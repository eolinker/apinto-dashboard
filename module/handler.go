package apinto_module

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
)

const (
	ApintoModuleName  = "apinto:module:name"
	ApintoRuntimeKeys = "apinto-runtime-keys"
)

func ModuleNameHandler(name string) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		ginCtx.Set(ApintoModuleName, name)
	}
}

func AddKey(ginCtx *gin.Context, key string, value interface{}) {
	ginCtx.Set(key, value)
	v := map[string]interface{}{
		key: value,
	}
	responseKeyData, err := json.Marshal(v)
	if err != nil {
		return
	}
	ginCtx.Writer.Header().Add(ApintoRuntimeKeys, string(responseKeyData))
}
