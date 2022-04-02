package profession

import (
	"fmt"
	apinto_dashboard "github.com/eolinker/apinto-dashboard"
	"net/http"
)

type Profession struct {
	views *apinto_dashboard.ModuleViewEmpty
}

func (p *Profession) Lookup(r *http.Request) (view string, data interface{}, err error) {
	return p.views.Lookup(r)
}

func (p *Profession) ServeHTTP(w http.ResponseWriter, req *http.Request) {

}

func NewProfession(name string) *Profession {
	views := map[string]string{
		"list":"profession_list",
		"create":"profession_create",
		"edit":"profession_edit",
	}
	return &Profession{
		views:	apinto_dashboard.NewViewModuleEmpty(fmt.Sprint("/", name, "/"), views, "list"),
	}
}