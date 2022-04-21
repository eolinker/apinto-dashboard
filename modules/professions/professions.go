package professions

import (
	"fmt"
	apinto_dashboard "github.com/eolinker/apinto-dashboard"
	"github.com/eolinker/apinto-dashboard/internal/apinto"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type ListHeader struct {
	Title  map[apinto_dashboard.ZoneName][]string
	Fields []string
}
type Profession struct {
	*apinto_dashboard.ModuleViewFinder
	*httprouter.Router
	ModuleName     string
	ProfessionName string

	header *ListHeader
}

func (p *Profession) Lookup(r *http.Request) (view string, data interface{}, has bool) {
	name, has := p.ModuleViewFinder.Lookup(r)
	if has {
		switch name {
		case "profession_list":
			return name, p.header, true
		case "profession_edit":
			workerName := r.URL.Query().Get("name")
			return name, workerName, true
		}
		return name, nil, true
	}
	return "", nil, false
}

func (p *Profession) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	p.Router.ServeHTTP(w, req)
}

func NewProfession(name string, profession string, titles map[apinto_dashboard.ZoneName][]string, fields []string, viewFinder *apinto_dashboard.ModuleViewFinder) *Profession {

	if viewFinder == nil {
		views := map[string]string{
			"list":   "profession_list",
			"create": "profession_create",
			"edit":   "profession_edit",
		}
		viewFinder = apinto_dashboard.NewViewModuleEmpty(fmt.Sprint("/", name, "/"), views, "list")
	}
	p := &Profession{
		ModuleName:       name,
		ProfessionName:   profession,
		ModuleViewFinder: viewFinder,
		header: &ListHeader{
			Title:  titles,
			Fields: fields,
		},
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
	r.GET(fmt.Sprintf("/api/%s/drivers", p.ModuleName), func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		data, code, err := apinto.Client().Drivers(p.ProfessionName)
		if err != nil {
			apinto.WriteResult(w, 500, []byte(err.Error()))
			return
		}
		apinto.WriteResult(w, code, data)
	})

	r.GET(fmt.Sprintf("/api/%s/render/:driver", p.ModuleName), func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		driver := params.ByName("driver")
		data, code, err := apinto.Client().Render(p.ProfessionName, driver)
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
