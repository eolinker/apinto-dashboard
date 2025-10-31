package builder

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type noRouterPrefix struct {
	prefixFrontendAssets map[string]struct{}
	PrefixApis           map[string]struct{}
}

func newNoRouterPrefix() *noRouterPrefix {
	return &noRouterPrefix{
		prefixFrontendAssets: make(map[string]struct{}),
		PrefixApis:           make(map[string]struct{}),
	}
}

var (
	_404Json = []byte(`{"code":404,"msg":"404 page not found"}`)
)

func (p *noRouterPrefix) addFrontendAssets(path string) {
	pre := strings.Split(strings.Trim(path, "/"), "/")[0]
	if len(pre) > 0 {
		p.prefixFrontendAssets[pre] = struct{}{}
	}
}
func (p *noRouterPrefix) addApi(path string) {
	pre := strings.Split(strings.Trim(path, "/"), "/")[0]
	if len(pre) > 0 {
		p.PrefixApis[pre] = struct{}{}
	}
}
func (p *noRouterPrefix) noRouterHandleFunc(ctx *gin.Context) {

	if ctx.Request.Method != http.MethodGet && ctx.Request.Method != http.MethodHead {
		ctx.Data(http.StatusNotFound, "application/json", _404Json)
		return
	}

	prefix := strings.Split(strings.Trim(ctx.Request.RequestURI, "/"), "/")[0]
	if len(prefix) == 0 {
		indexHtmlHandle(ctx)
		return
	}
	if _, has := p.prefixFrontendAssets[prefix]; has {
		ctx.Data(http.StatusNotFound, "application/text", []byte("404 page not found"))
		return
	}
	if _, has := p.PrefixApis[prefix]; has {
		ctx.Data(http.StatusNotFound, "application/json", _404Json)
		return
	}
	indexHtmlHandle(ctx)
	return

}
