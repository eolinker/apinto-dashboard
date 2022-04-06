package activity_log

import (
	"net/http"
)

type ActivityLog struct {
	name string
}

func NewActivityLog(name string) *ActivityLog {
	return &ActivityLog{name: name}
}

func (a *ActivityLog) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

func (a *ActivityLog) Lookup(r *http.Request) (view string, data interface{}, has bool) {
	return a.name,nil,true
}

