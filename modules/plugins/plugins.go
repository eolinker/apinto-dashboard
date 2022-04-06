package plugins

import (
	"net/http"
)

type Plugins struct {
	name string
}

func NewPlugins(name string) *Plugins {
	return &Plugins{name: name}
}

func (p *Plugins) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 这里处理 插件api
}

func (p *Plugins) Lookup(r *http.Request) (view string, data interface{}, has bool) {
	return p.name,"",true
}
