package apinto_dashboard

import (
	"errors"
	"net/http"
	"strings"
)

type IModule interface {
	http.Handler
	Lookup(r *http.Request)(view string,data interface{},err error)
}


type ModuleViewEmpty struct {
	views map[string]string
	prefix string
	defaultName string
}



func NewViewModuleEmpty(prefix string, views map[string]string, defaultName string) *ModuleViewEmpty {
	 prefix = strings.TrimSuffix(prefix,"/")
	return &ModuleViewEmpty{
		views: views,
		prefix: prefix,
		defaultName: defaultName,
	}
}

func (v *ModuleViewEmpty) Lookup(r *http.Request) (view string, data interface{}, err error) {
	path := r.URL.Path
	name:= strings.TrimPrefix(path,v.prefix)
	if name=="" || name =="/"{
		name = v.defaultName
	}
	 name = strings.TrimPrefix(name,"/")
	vn,has:=v.views[name]
	if has{
		return vn,nil,nil
	}
	return "",nil,errors.New("404 page not found")
}

