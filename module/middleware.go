package apinto_module

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GenHandleMiddleware(handle func(ctx context.Context, request *MiddlewareRequest, writer MiddlewareResponseWriter)) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		ginCtx.FullPath()
		request, err := Read(ginCtx.Request)
		if err != nil {
			ginCtx.AbortWithError(http.StatusServiceUnavailable, err)
			return
		}

		writer := new(MiddlewareResponse)
		handle(ginCtx, request, writer)
		ginCtx.JSON(http.StatusOK, writer)

	}
}
