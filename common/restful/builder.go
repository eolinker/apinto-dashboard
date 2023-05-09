package restful

import (
	"net/http"
	"net/url"
)

type requestCore struct {
	*builderCore
	header http.Header
	query  url.Values
}

type builderCore struct {
	Config
	path   PathGen
	method string
}

func newBuilderCore(config Config, method string, path string) builderCore {
	return builderCore{Config: config, path: createPath(path), method: method}
}

func (b *builderCore) build() *requestCore {

	return &requestCore{builderCore: b, header: b.header.Clone(), query: cloneValues(b.query)}

}
func (b *builderCore) Reset(c Config) {
	b.Config = c
}
