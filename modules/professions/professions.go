package professions

import (
	"fmt"
	apinto_dashboard "github.com/eolinker/apinto-dashboard"
	"net/http"
)

type Profession struct {
	views *apinto_dashboard.ModuleViewFinder
	professionName string
}

func (p *Profession) Lookup(r *http.Request) (view string, data interface{}, has bool) {
	name, has := p.views.Lookup(r)
	if has{
		return name,nil,true
	}
	return "",nil,false
}

func (p *Profession) ServeHTTP(w http.ResponseWriter, req *http.Request) {

}

func NewProfession(name string,profession string) *Profession {
	views := map[string]string{
		"list":"profession_list",
		"create":"profession_create",
		"edit":"profession_edit",
	}
	return &Profession{
		professionName:profession,
		views:	apinto_dashboard.NewViewModuleEmpty(fmt.Sprint("/", name, "/"), views, "list"),
	}
}