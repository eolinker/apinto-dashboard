package apinto_module

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	ApintoModuleName  = "apinto:module:name"
	ApintoRuntimeKeys = "apinto-runtime-keys"
	ApintoProxyEvent  = "apinto-proxy-event"
)

func ModuleNameHandler(name string) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		ginCtx.Set(ApintoModuleName, name)
	}
}

func RebuildKeys(ginCtx *gin.Context) {
	data := ginCtx.GetHeader(ApintoRuntimeKeys)
	keys := make(map[string]any)
	err := json.Unmarshal([]byte(data), &keys)
	if err != nil {
		return
	}
	for k, v := range keys {
		ginCtx.Set(k, v)
	}
	ginCtx.Next()
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
func ReSetHeaderFromProxyResponse(responseHeader http.Header, ginCtx *gin.Context) {

	for k, vs := range responseHeader {
		switch strings.ToLower(k) {
		case ApintoProxyEvent:
			for _, v := range vs {
				if len(v) > 0 {
					eventObjs := make(map[string]any)
					if errError := json.Unmarshal([]byte(v), &eventObjs); errError == nil {
						for event, ev := range eventObjs {
							DoEvent(event, ev)
						}
					}
				}
			}
		case ApintoRuntimeKeys:
			for _, v := range vs {
				keys := make(map[string]any)
				if errError := json.Unmarshal([]byte(v), &keys); errError == nil {
					for key, ev := range keys {
						ginCtx.Set(key, ev)
					}
				}
			}
		default:
			for _, v := range vs {
				ginCtx.Writer.Header().Add(k, v)
			}

		}
	}
}
func AddEvent(ginCtx *gin.Context, event string, value any) {
	v := map[string]any{
		event: value,
	}
	eventData, err := json.Marshal(v)
	if err != nil {
		return
	}
	ginCtx.Writer.Header().Add(ApintoProxyEvent, string(eventData))
}
func CopyKeysToHeader(keys map[string]any, req *http.Request) {
	if keys == nil {
		return
	}
	apintoKeysData, err := json.Marshal(keys)
	if err != nil {
		return
	}
	req.Header.Set(ApintoRuntimeKeys, string(apintoKeysData))
}
