package pm3

import (
	"fmt"
	"net/http"
)

type ModuleTool struct {
	id   string
	name string
}

func NewModuleTool(id, name string) *ModuleTool {
	return &ModuleTool{id: id, name: name}
}
func (m *ModuleTool) Name() string {
	return m.name
}

func (m *ModuleTool) InitAccess(apis []Api) {
	for i := range apis {
		if apis[i].Authority == Anonymous || apis[i].Authority == Public {
			//fmt.Printf("InitAccess %s:%s %s\n", apis[i].Authority.String(), apis[i].Method, apis[i].Path)
			continue
		}
		if apis[i].Access == "" {
			switch apis[i].Method {
			case http.MethodGet:
				apis[i].Access = fmt.Sprintf("%s.%s.view", m.id, m.name)
			case http.MethodPost, http.MethodDelete, http.MethodPut:
				apis[i].Access = fmt.Sprintf("%s.%s.edit", m.id, m.name)
			}
		}
		if apis[i].Authority.String() == "unset" {
			switch apis[i].Method {
			case http.MethodGet:
				apis[i].Authority = Internal
			case http.MethodPost, http.MethodPut, http.MethodDelete:
				apis[i].Authority = Private

			}
		}
		//fmt.Printf("InitAccess %s:%s %s=>%s\n", apis[i].Authority, apis[i].Method, apis[i].Path, apis[i].Access)
	}
}
