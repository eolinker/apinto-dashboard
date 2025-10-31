package apinto_module

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/pm3"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
	"strings"
)

const (
	filePathParamName = "filepath"
)

func StaticFileDir(prefix string, fileSystem http.FileSystem) pm3.FrontendAsset {

	handler := func(ginCtx *gin.Context) {
		filePath := ginCtx.Param(filePathParamName)
		ginCtx.FileFromFS(filePath, fileSystem)
	}
	if !strings.HasSuffix(prefix, "/") {
		prefix = fmt.Sprint(prefix, "/")
	}
	return pm3.FrontendAsset{
		Path:        prefix,
		HandlerFunc: handler,
	}
}

func StaticRouter(prefix string) string {
	if i := strings.Index(prefix, "*"); i > 0 {
		prefix = prefix[:i]
	}
	if !strings.HasSuffix(prefix, "/") {
		return prefix
	}
	return path.Join(prefix, fmt.Sprint("*", filePathParamName))
}
