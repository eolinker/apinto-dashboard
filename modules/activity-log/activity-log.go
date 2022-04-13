package activity_log

import (
	"encoding/json"
	"fmt"
	apinto_dashboard "github.com/eolinker/apinto-dashboard"
	"github.com/eolinker/apinto-dashboard/internal/apinto"
	"github.com/eolinker/apinto-dashboard/modules/activity-log/db"
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

func NewActivityLog(name string) (*ActivityLog, error) {
	err := db.InitDB()
	if err != nil {
		return nil, err
	}

	views := map[string]string{
		"list": "activity_log",
	}

	activityLog := &ActivityLog{
		ModuleViewFinder: apinto_dashboard.NewViewModuleEmpty(fmt.Sprint("/", name, "/"), views, "list"),
		ModuleName:       name,
		header: &professions.ListHeader{
			Title: map[apinto_dashboard.ZoneName][]string{
				apinto_dashboard.ZhCn: {"操作时间", "用户", "内容"},
				apinto_dashboard.EnUs: {"Time", "User", "Content"},
			},
			Fields: []string{"time", "user", "content"},
		},
	}
	activityLog.createRouter()

	return activityLog, nil
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
			"user":    "admin",
			"content": "操作:创建 操作对象:demoRouter",
			"time":    time.Now().Format("2006-01-02 15:04:05"),
		}}
		data, _ := json.Marshal(fakeData)

		apinto.WriteResult(w, 200, data)
	})

	a.Router = r
}
