package monitors

import (
	"net/http"
)

type Monitor struct {
	name string
}

func NewMonitor(name string) *Monitor {
	return &Monitor{name: name}
}

func (m *Monitor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	
}

func (m *Monitor) Lookup(r *http.Request) (view string, data interface{}, has bool) {
	return m.name,nil,true
}
