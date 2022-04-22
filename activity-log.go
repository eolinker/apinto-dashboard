package apinto_dashboard

import "log"

type ActivityLogAddHandler interface {
	Add(user, content, operation, object string, args []*Arg) error
}

type ActivityLogGetHandler interface {
	GetLogList(offset, limit int, user, operation, object string, startUnix, endUnix int64) ([]*LogEntity, int64, error)
}
type LogEntity struct {
	Time      string `json:"time"`
	User      string `json:"user"`
	Operation string `json:"operation"`
	Object    string `json:"object"`
	Content   string `json:"content"`
	Args      []*Arg `json:"args"`
}
type Arg struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

var (
	activityLogHandler ActivityLogAddHandler
)

func SetActivityLogAddHandler(h ActivityLogAddHandler) {
	activityLogHandler = h
}
func AddActivityLog(user, operation, object, content string, args []*Arg) {
	if activityLogHandler != nil {
		err := activityLogHandler.Add(user, operation, object, content, args)
		if err != nil {
			log.Println("[ERR] add log:", err)
		}
	}
}
