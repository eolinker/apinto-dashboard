package routers

import (
	apinto_dashboard "github.com/eolinker/apinto-dashboard"
	"net/http"
)

type Routers struct {
	views *apinto_dashboard.ModuleViewFinder
}

func (p *Routers) Lookup(r *http.Request) (view string, data interface{}, has bool) {
	name, has := p.views.Lookup(r)
	if has{
		return name,nil,true
	}
	return "",nil,false
}

func (p *Routers) ServeHTTP(w http.ResponseWriter, req *http.Request) {

}

func NewRouters() *Routers {
	views := map[string]string{
		"list":"router_list",
		"create":"router_create",
		"edit":"router_edit",
	}
	return &Routers{
		views:	apinto_dashboard.NewViewModuleEmpty("/routers/", views, "list"),
	}
}