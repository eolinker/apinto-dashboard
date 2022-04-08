package professions

import (
	"fmt"
	apinto_dashboard "github.com/eolinker/apinto-dashboard"
	"github.com/eolinker/apinto-dashboard/internal/apinto"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Profession struct {
	views          *apinto_dashboard.ModuleViewFinder
	ModuleName     string
	ProfessionName string
	Router         *httprouter.Router
}

func (p *Profession) Lookup(r *http.Request) (view string, data interface{}, has bool) {
	name, has := p.views.Lookup(r)
	if has {
		return name, nil, true
	}
	return "", nil, false
}

func (p *Profession) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	p.Router.ServeHTTP(w, req)
}

func NewProfession(name string, profession string) *Profession {
	views := map[string]string{
		"list":   "profession_list",
		"create": "profession_create",
		"edit":   "profession_edit",
	}
	p := &Profession{
		ModuleName:     name,
		ProfessionName: profession,
		views:          apinto_dashboard.NewViewModuleEmpty(fmt.Sprint("/", name, "/"), views, "list"),
	}
	p.createRouter()
	return p
}

func (p *Profession) createRouter() {
	r := httprouter.New()

	// List
	r.GET(fmt.Sprintf("/api/%s/", p.ModuleName), func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		data, code, err := apinto.Client().List(p.ProfessionName)
		if err != nil {
			apinto.WriteResult(w, 500, []byte(err.Error()))
			return
		}
		apinto.WriteResult(w, code, data)
	})

	// Get
	r.GET(fmt.Sprintf("/api/%s/:name", p.ModuleName), func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		name := params.ByName("name")
		data, code, err := apinto.Client().Get(p.ProfessionName, name)
		if err != nil {
			apinto.WriteResult(w, 500, []byte(err.Error()))
			return
		}
		apinto.WriteResult(w, code, data)
	})

	// Create
	r.POST(fmt.Sprintf("/api/%s/", p.ModuleName), func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		data, err := apinto.ReadBody(r.Body)
		if err != nil {
			apinto.WriteResult(w, 500, []byte(err.Error()))
			return
		}
		data, code, err := apinto.Client().Create(p.ProfessionName, data)
		if err != nil {
			apinto.WriteResult(w, 500, []byte(err.Error()))
			return
		}
		apinto.WriteResult(w, code, data)
	})

	// Update
	r.PUT(fmt.Sprintf("/api/%s/:name", p.ModuleName), func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		name := params.ByName("name")
		data, err := apinto.ReadBody(r.Body)
		if err != nil {
			apinto.WriteResult(w, 500, []byte(err.Error()))
			return
		}
		data, code, err := apinto.Client().Update(p.ProfessionName, name, data)
		if err != nil {
			apinto.WriteResult(w, 500, []byte(err.Error()))
			return
		}
		apinto.WriteResult(w, code, data)
	})

	// Delete
	r.DELETE(fmt.Sprintf("/api/%s/:name", p.ModuleName), func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		name := params.ByName("name")
		data, code, err := apinto.Client().Delete(p.ProfessionName, name)
		if err != nil {
			apinto.WriteResult(w, 500, []byte(err.Error()))
			return
		}
		apinto.WriteResult(w, code, data)
	})
	p.Router = r
}
