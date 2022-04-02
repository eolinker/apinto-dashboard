package apinto_dashboard

import (
	"net/http"
	"strings"
)

type IModule interface {
	http.Handler
	Lookup(r *http.Request) (view string, data interface{}, has bool)
}


type ModuleViewFinder struct {
	views map[string]string
	prefix string
	defaultName string
}



func NewViewModuleEmpty(prefix string, views map[string]string, defaultName string) *ModuleViewFinder {
	 prefix = strings.TrimSuffix(prefix,"/")
	return &ModuleViewFinder{
		views: views,
		prefix: prefix,
		defaultName: defaultName,
	}
}

func (v *ModuleViewFinder) Lookup(r *http.Request) (view string,  has bool) {
	path := r.URL.Path
	name:= strings.TrimPrefix(path,v.prefix)
	if name=="" || name =="/"{
		name = v.defaultName
	}
	 name = strings.TrimPrefix(name,"/")
	vn,has:=v.views[name]

	return vn,has
}

