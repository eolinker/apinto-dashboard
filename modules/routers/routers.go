package routers

import (
	"fmt"
	"net/http"

	apinto_dashboard "github.com/eolinker/apinto-dashboard"
	"github.com/eolinker/apinto-dashboard/internal/apinto"
	"github.com/eolinker/apinto-dashboard/modules/professions"
	"github.com/julienschmidt/httprouter"
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
			return name, routerName, true
		case "router_create":
			return name, p.ProfessionName, true
		}

		return name, nil, true
	}

	return "", nil, false
}

func NewRouters(name string) *Routers {
	views := map[string]string{
		"list":   "router_list",
		"create": "router_create",
		"edit":   "router_edit",
	}
	professionsHandler := professions.NewProfession(name, "router",
		nil, nil,
		apinto_dashboard.NewViewModuleEmpty("/routers/", views, "list"))
	r := &Routers{
		Profession: professionsHandler,
		header: &professions.ListHeader{
			Title: map[apinto_dashboard.ZoneName][]string{
				apinto_dashboard.ZhCn: {"路由名", "驱动", "域名", "端口", "服务", "创建时间", "更新时间"},
				apinto_dashboard.EnUs: {"Name", "Driver", "Host", "Listen", "Service", "Create", "Update"},
			},
			Fields: []string{"name", "driver", "host", "listen", "target", "create", "update"},
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
		rdata, code, err := apinto.Client().Patch(p.ProfessionName, name, data)
		if err != nil {
			apinto.WriteResult(w, 500, []byte(err.Error()))
			return
		}
		apinto.WriteResult(w, code, rdata)
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
