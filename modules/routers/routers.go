package routers

import (
	"fmt"
	apinto_dashboard "github.com/eolinker/apinto-dashboard"
	"github.com/eolinker/apinto-dashboard/internal/apinto"
	"github.com/eolinker/apinto-dashboard/modules/professions"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Routers struct {
	*professions.Profession
	header *professions.ListHeader
}

func (p *Routers) Lookup(r *http.Request) (view string, data interface{}, has bool) {
	name, has := p.ModuleViewFinder.Lookup(r)
	if has {
		switch name {
		case "router_list":
			return name, p.header, true
		case "router_edit":
			routerName := r.URL.Query().Get("name")
			driver := r.URL.Query().Get("driver")
			d := map[string]string{
				"profession": routerName,
				"driver":     driver,
			}
			return name, d, true
		}
		return name, nil, true
	}

	return "", nil, false
}

func (p *Routers) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	p.Router.ServeHTTP(w, req)
}

func NewRouters() *Routers {
	views := map[string]string{
		"list":   "router_list",
		"create": "router_create",
		"edit":   "router_edit",
	}
	professionsHandler := professions.NewProfession("routers", "router",
		nil, nil,
		apinto_dashboard.NewViewModuleEmpty("/routers/", views, "list"))
	r := &Routers{
		Profession: professionsHandler,
		header: &professions.ListHeader{
			Title: map[apinto_dashboard.ZoneName][]string{
				apinto_dashboard.ZhCn: {"路由名", "驱动", "域名", "端口", "服务", "状态", "创建时间", "更新时间"},
				apinto_dashboard.EnUs: {"Name", "Driver", "Host", "Listen", "Service", "Status", "Create", "Update"},
			},
			Fields: []string{"name", "driver", "host", "listen", "service", "status", "create", "update"},
		},
	}
	r.expandRouter()
	return r
}

func (p *Routers) expandRouter() {
	p.Router.PATCH(fmt.Sprintf("/api/%s/:name", p.ModuleName), func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		name := params.ByName("name")
		data, err := apinto.ReadBody(r.Body)
		if err != nil {
			apinto.WriteResult(w, 500, []byte(err.Error()))
			return
		}
		data, code, err := apinto.Client().Patch(p.ProfessionName, name, data)
		if err != nil {
			apinto.WriteResult(w, 500, []byte(err.Error()))
			return
		}
		apinto.WriteResult(w, code, data)
	})

	// PatchPath
	p.Router.PATCH(fmt.Sprintf("/api/%s/:name/*path", p.ModuleName), func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		name := params.ByName("name")
		path := params.ByName("path")
		data, err := apinto.ReadBody(r.Body)
		if err != nil {
			apinto.WriteResult(w, 500, []byte(err.Error()))
			return
		}
		data, code, err := apinto.Client().PatchPath(p.ProfessionName, name, path, data)
		if err != nil {
			apinto.WriteResult(w, 500, []byte(err.Error()))
			return
		}
		apinto.WriteResult(w, code, data)
	})
}
