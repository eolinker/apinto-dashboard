package plugins

import (
	"fmt"
	apinto_dashboard "github.com/eolinker/apinto-dashboard"
	"github.com/eolinker/apinto-dashboard/internal/apinto"
	"github.com/eolinker/apinto-dashboard/modules/professions"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Plugins struct {
	header *professions.ListHeader
	*httprouter.Router
	name           string
	ProfessionName string
	workerName     string
}

func NewPlugins(name string) *Plugins {

	p := &Plugins{
		name:           name,
		ProfessionName: "setting",
		workerName:     "plugin",
		header: &professions.ListHeader{
			Title: map[apinto_dashboard.ZoneName][]string{
				apinto_dashboard.ZhCn: {"名称", "扩展ID", "类型", "状态"},
				apinto_dashboard.EnUs: {"Name", "ID", "Type", "Status"},
			},
			Fields: []string{"name", "id", "type", "status"},
		},
	}
	p.createRouter()
	return p
}

func (p *Plugins) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 这里处理 插件api
	p.Router.ServeHTTP(w, r)
}

func (p *Plugins) Lookup(r *http.Request) (view string, data interface{}, has bool) {
	return p.name, p.header, true
}

func (p *Plugins) createRouter() {
	r := httprouter.New()
	// List
	r.GET(fmt.Sprintf("/api/%s/", p.name), func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		data, code, err := apinto.Client().Get(p.ProfessionName, p.workerName)
		if err != nil {
			apinto.WriteResult(w, 500, []byte(err.Error()))
			return
		}
		apinto.WriteResult(w, code, data)
	})
	// Update
	r.PUT(fmt.Sprintf("/api/%s/", p.name), func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		data, err := apinto.ReadBody(r.Body)
		if err != nil {
			apinto.WriteResult(w, 500, []byte(err.Error()))
			return
		}
		data, code, err := apinto.Client().Update(p.ProfessionName, p.workerName, data)
		if err != nil {
			apinto.WriteResult(w, 500, []byte(err.Error()))
			return
		}
		apinto.WriteResult(w, code, data)
	})
	p.Router = r
}
