package activity_log

import (
	"fmt"
	apinto_dashboard "github.com/eolinker/apinto-dashboard"
	"github.com/eolinker/apinto-dashboard/internal/apinto"
	"github.com/eolinker/apinto-dashboard/modules/activity-log/database"
	"github.com/eolinker/apinto-dashboard/modules/activity-log/database/module"
	"github.com/eolinker/apinto-dashboard/modules/professions"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type ActivityLog struct {
	*apinto_dashboard.ModuleViewFinder
	*httprouter.Router
	ModuleName string
	header     *professions.ListHeader
}

func NewActivityLog(name string) (*ActivityLog, error) {
	err := database.InitDB()
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
				apinto_dashboard.ZhCn: {"操作时间", "用户", "操作类型", "操作对象", "内容"},
				apinto_dashboard.EnUs: {"Time", "User", "operation", "object", "Content"},
			},
			Fields: []string{"time", "user", "operation", "object", "content"},
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
		//从sqlite读取数据
		offsetStr := r.URL.Query().Get("offset")
		limitStr := r.URL.Query().Get("limit")
		user := r.URL.Query().Get("user")
		operation := r.URL.Query().Get("operation")
		object := r.URL.Query().Get("object")
		startUnixStr := r.URL.Query().Get("startUnix")
		endUnixStr := r.URL.Query().Get("endUnix")

		offset, _ := strconv.Atoi(offsetStr)
		limit, _ := strconv.Atoi(limitStr)
		startUnix, _ := strconv.ParseInt(startUnixStr, 10, 64)
		endUnix, _ := strconv.ParseInt(endUnixStr, 10, 64)

		data, err := module.GetLogList(offset, limit, user, operation, object, startUnix, endUnix)
		if err != nil {
			apinto.WriteResult(w, 500, []byte(err.Error()))
			return
		}

		apinto.WriteResult(w, 200, data)
	})

	a.Router = r
}
