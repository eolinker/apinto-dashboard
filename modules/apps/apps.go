package apps

import (
	"fmt"
	"net/http"

	apinto_dashboard "github.com/eolinker/apinto-dashboard"
	"github.com/eolinker/apinto-dashboard/internal/apinto"
	"github.com/eolinker/apinto-dashboard/modules/professions"
	"github.com/julienschmidt/httprouter"
)

type Apps struct {
	*professions.Profession
	header *professions.ListHeader
}

func (p *Apps) Lookup(r *http.Request) (view string, data interface{}, has bool) {
	name, has := p.ModuleViewFinder.Lookup(r)
	if has {
		switch name {
		case "app_list":
			return name, p.header, true
		case "profession_edit":
			routerName := r.URL.Query().Get("name")
			return name, routerName, true
		case "profession_create":
			return name, p.ProfessionName, true
		}

		return name, nil, true
	}

	return "", nil, false
}

func NewApps(name string) *Apps {
	views := map[string]string{
		"list":   "app_list",
		"create": "profession_create",
		"edit":   "profession_edit",
	}
	professionsHandler := professions.NewProfession(name, "app", nil, nil,
		apinto_dashboard.NewViewModuleEmpty("/apps/", views, "list"))
	r := &Apps{
		Profession: professionsHandler,
		header: &professions.ListHeader{
			Title: map[apinto_dashboard.ZoneName][]string{
				apinto_dashboard.ZhCn: {"应用名称", "驱动", "创建时间", "更新时间"},
				apinto_dashboard.EnUs: {"Name", "Driver", "Create", "Update"},
			},
			Fields: []string{"name", "driver", "create", "update"},
		},
	}
	r.expandApp()
	return r
}

func (p *Apps) expandApp() {
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
