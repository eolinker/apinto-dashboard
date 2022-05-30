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
		userName := apinto_dashboard.MustUsername(r)

		logContent := "编辑全局插件失败"
		logArgs := make([]*apinto_dashboard.Arg, 0, 6)
		logArgs = append(logArgs,
			&apinto_dashboard.Arg{Key: "user", Value: userName},
			&apinto_dashboard.Arg{Key: "profession", Value: p.ProfessionName},
			&apinto_dashboard.Arg{Key: "url", Value: r.URL.String()},
		)

		defer func() {
			apinto_dashboard.AddActivityLog(r, userName, apinto_dashboard.OPT_UPDATE, "全局插件", logContent, logArgs)
		}()

		rData, err := apinto.ReadBody(r.Body)
		if err != nil {
			logContent = "编辑全局插件失败:Body读取失败"
			logArgs = append(logArgs,
				&apinto_dashboard.Arg{Key: "error", Value: err.Error()},
				&apinto_dashboard.Arg{Key: "err_from", Value: "dashboard"},
			)

			apinto.WriteResult(w, 500, []byte(err.Error()))
			return
		}

		//更新全局插件
		data, code, err := apinto.Client().Update(p.ProfessionName, p.workerName, rData)
		if err != nil {
			logArgs = append(logArgs,
				&apinto_dashboard.Arg{Key: "request_body", Value: string(rData)},
				&apinto_dashboard.Arg{Key: "error", Value: err.Error()},
				&apinto_dashboard.Arg{Key: "err_from", Value: "dashboard"},
			)

			apinto.WriteResult(w, 500, []byte(err.Error()))
		} else if code != 200 {
			logArgs = append(logArgs,
				&apinto_dashboard.Arg{Key: "request_body", Value: string(rData)},
				&apinto_dashboard.Arg{Key: "error", Value: string(data)},
				&apinto_dashboard.Arg{Key: "err_from", Value: "apinto"},
			)

			apinto.WriteResult(w, code, data)
		} else {
			logContent = "编辑全局插件成功"
			logArgs = append(logArgs,
				&apinto_dashboard.Arg{Key: "request_body", Value: string(rData)},
			)

			apinto.WriteResult(w, code, data)
		}

	})
	p.Router = r
}
