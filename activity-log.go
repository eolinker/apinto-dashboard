package apinto_dashboard

import (
	"log"
	"net/http"
	"strings"
)

type ActivityLogAddHandler interface {
	Add(user, ip, content, operation, target string, args []*Arg) error
}

type ActivityLogGetHandler interface {
	GetLogList(offset, limit int, user, operation, target string, startUnix, endUnix int64) ([]*LogEntity, int64, error)
}
type LogEntity struct {
	Time      string `json:"time"`
	User      string `json:"user"`
	IP        string `json:"ip"`
	Operation string `json:"operation"`
	Target    string `json:"target"`
	Content   string `json:"content"`
	Args      []*Arg `json:"args"`
}
type Arg struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

var (
	activityLogHandler ActivityLogAddHandler
	isFileterForwarded bool
)

func SetActivityLogAddHandler(h ActivityLogAddHandler, fileterForwarded bool) {
	activityLogHandler = h
	isFileterForwarded = fileterForwarded
}
func AddActivityLog(r *http.Request, user, operation, target, content string, args []*Arg) {
	if activityLogHandler != nil {
		err := activityLogHandler.Add(user, getIP(r), operation, target, content, args)
		if err != nil {
			log.Println("[ERR] add log:", err)
		}
	}
}

func getIP(r *http.Request) string {
	if !isFileterForwarded {
		forwarded := r.Header.Get("x-forwarded-for")
		if len(forwarded) > 0 {
			if i := strings.Index(forwarded, ","); i > 0 {
				return forwarded[:i]
			}
			return forwarded
		}
	}

	remoteIP := r.RemoteAddr
	idx := strings.LastIndex(remoteIP, ":")
	if idx > 0 {
		return remoteIP[:idx]
	}
	return remoteIP
}
