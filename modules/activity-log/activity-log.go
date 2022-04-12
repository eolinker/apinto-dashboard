package activity_log

import (
	"encoding/json"
	"fmt"
	apinto_dashboard "github.com/eolinker/apinto-dashboard"
	"github.com/eolinker/apinto-dashboard/internal/apinto"
	"github.com/eolinker/apinto-dashboard/modules/professions"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

type ActivityLog struct {
	*apinto_dashboard.ModuleViewFinder
	*httprouter.Router
	ModuleName string

	header *professions.ListHeader
}

func NewActivityLog(name string) *ActivityLog {
	views := map[string]string{
		"list": "activity_log",
	}

	activityLog := &ActivityLog{
		ModuleViewFinder: apinto_dashboard.NewViewModuleEmpty(fmt.Sprint("/", name, "/"), views, "list"),
		ModuleName:       name,
		header: &professions.ListHeader{
			Title: map[apinto_dashboard.ZoneName][]string{
				apinto_dashboard.ZhCn: {"序号", "用户", "操作", "内容", "操作时间"},
				apinto_dashboard.EnUs: {"Order", "User", "Operation", "Content", "Time"},
			},
			Fields: []string{"order", "user", "operation", "content", "time"},
		},
	}
	activityLog.createRouter()

	return activityLog
}

func (a *ActivityLog) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.Router.ServeHTTP(w, r)
}

func (a *ActivityLog) Lookup(r *http.Request) (view string, data interface{}, has bool) {
	name, has := a.ModuleViewFinder.Lookup(r)
	if has {
		switch name {
		case "activity_log":
			return name, a.header, true
		}
		return name, nil, true
	}

	return "", nil, false
}

func (a *ActivityLog) createRouter() {
	r := httprouter.New()

	// List
	r.GET(fmt.Sprintf("/api/%s/", a.ModuleName), func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		//TODO  从sqlit读取数据

		//if err != nil {
		//	apinto.WriteResult(w, 500, []byte(err.Error()))
		//	return
		//}
		//apinto.WriteResult(w, code, data)
		fakeData := []map[string]interface{}{{
			"order":     1,
			"user":      "admin",
			"operation": "Create",
			"content":   "{\"name\":\"demo\",\"driver\":\"http\",\"desc\":\"http\",\"listen\":8080,\"rules\":[{\"location\":\"/Web/Test/params/print\"}],\"target\":\"demo@service\"}",
			"time":      time.Now().String(),
		}}
		data, _ := json.Marshal(fakeData)

		apinto.WriteResult(w, 200, data)
	})

	a.Router = r
}
