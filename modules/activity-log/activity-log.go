package activity_log

import (
	"encoding/json"
	"fmt"
	apinto "github.com/eolinker/apinto-dashboard"
	response "github.com/eolinker/apinto-dashboard/internal/apinto"
	"github.com/eolinker/apinto-dashboard/modules/professions"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type ActivityLog struct {
	*apinto.ModuleViewFinder
	*httprouter.Router
	ModuleName string
	header     *professions.ListHeader
	dao        apinto.ActivityLogGetHandler
}

func NewActivityLog(name string, dao apinto.ActivityLogGetHandler) (*ActivityLog, error) {

	views := map[string]string{
		"list": "activity_log",
	}

	activityLog := &ActivityLog{
		ModuleViewFinder: apinto.NewViewModuleEmpty(fmt.Sprint("/", name, "/"), views, "list"),
		ModuleName:       name,
		header: &professions.ListHeader{
			Title: map[apinto.ZoneName][]string{
				apinto.ZhCn: {"操作时间", "用户", "操作类型", "操作对象", "内容", "IP"},
				apinto.EnUs: {"Time", "User", "operation", "target", "Content", "IP"},
			},
			Fields: []string{"time", "user", "operation", "target", "content", "ip"},
		},
		dao: dao,
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
		target := r.URL.Query().Get("target")
		startUnixStr := r.URL.Query().Get("startUnix")
		endUnixStr := r.URL.Query().Get("endUnix")

		offset, _ := strconv.Atoi(offsetStr)
		limit, _ := strconv.Atoi(limitStr)
		startUnix, _ := strconv.ParseInt(startUnixStr, 10, 64)
		endUnix, _ := strconv.ParseInt(endUnixStr, 10, 64)

		data, err := a.getLogList(offset, limit, user, operation, target, startUnix, endUnix)
		if err != nil {
			response.WriteResult(w, 500, []byte(err.Error()))
			return
		}
		response.WriteResult(w, 200, data)
	})

	a.Router = r
}
func (a *ActivityLog) getLogList(offset, limit int, user, operation, target string, startUnix, endUnix int64) ([]byte, error) {
	list, total, err := a.dao.GetLogList(offset, limit, user, operation, target, startUnix, endUnix)
	if err != nil {
		return nil, err
	}
	m := make(map[string]interface{})
	m["list"] = list
	m["total_num"] = total

	data, _ := json.Marshal(m)
	return data, nil
}
