package render

import (
	"net/http"

	apinto_dashboard "github.com/eolinker/apinto-dashboard"
	"github.com/eolinker/apinto-dashboard/internal/apinto"
	"github.com/eolinker/apinto-dashboard/modules/professions"
	"github.com/julienschmidt/httprouter"
)

type Render struct {
	header *professions.ListHeader
	*httprouter.Router
	name           string
	ProfessionName string
}

func NewRender() *Render {
	p := &Render{
		ProfessionName: "setting",
		header: &professions.ListHeader{
			Title: map[apinto_dashboard.ZoneName][]string{
				apinto_dashboard.ZhCn: {"名称", "扩展ID", "状态"},
				apinto_dashboard.EnUs: {"Name", "ID", "Status"},
			},
			Fields: []string{"name", "id", "status"},
		},
	}
	p.createRouter()
	return p
}

func (p *Render) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 这里处理 插件api
	p.Router.ServeHTTP(w, r)
}

func (p *Render) Lookup(r *http.Request) (view string, data interface{}, has bool) {
	return p.name, p.header, true
}

func (p *Render) createRouter() {
	r := httprouter.New()
	// List
	r.GET("/setting/:name", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		data, code, err := apinto.Client().Get(p.ProfessionName, params.ByName("name"))
		if err != nil {
			apinto.WriteResult(w, 500, []byte(err.Error()))
			return
		}
		apinto.WriteResult(w, code, data)
	})

	p.Router = r
}
