package apinto_module

import "github.com/gin-gonic/gin"

type AccessReaderHandler struct {
	accessOfPath map[string][]string
}

func (a *AccessReaderHandler) ReadAccess(ctx *gin.Context) {
	access, has := a.accessOfPath[ctx.FullPath()]
	if has {
		ctx.Set("access", access)
	}
}

type AccessReader struct {
	handler gin.HandlerFunc
}

func (a *AccessReader) ReadAccess(ctx *gin.Context) {
	if a.handler == nil {
		return
	}
	a.handler(ctx)
}
