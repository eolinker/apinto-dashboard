package apinto_dashboard

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/internal/template"
	"net/http"
	"strings"
)

type IDashboardAccount interface {
}

type ViewHandler interface {
	Server(r *http.Request)(view string,data interface{},err error)
}

type Module struct {
	Path     string
	Icon     string
	Views    map[string]ViewHandler
	Handler  http.Handler
	Name     string
	I18nName map[string]string `json:"i18n_name"`
}

type Config struct {
	DefaultZone ZoneName
	Modules []*Module
	UserDetailsService IUserDetailsService
}
type DashboardService struct {
	defaultZone ZoneName
	userDetails IUserDetailsService
	server http.ServeMux
}

func (d *DashboardService) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 处理登陆
	d.server.ServeHTTP(w,req)
}

func Create(config *Config)(*DashboardService,error)  {
	if config.UserDetailsService == nil{
		return nil,ErrorUserDetailsServiceNeed
	}
	if config.DefaultZone == ""{
		config.DefaultZone = ZhCn
	}
	service := new(DashboardService)
	service.userDetails = config.UserDetailsService
	service.defaultZone = config.DefaultZone
	modules := make([]*ModuleItem,0,len(config.Modules))
	for _,m:=range config.Modules{
		modules = append(modules, &ModuleItem{
			Name:     m.Name,
			Icon:     m.Icon,
			I18nName: m.I18nName,
			Path:     m.Path,
		})
	}

	for _, m:=range config.Modules{

		prefix := fmt.Sprint("/", m.Name)

		for path,view:=range m.Views {
			if len(path)>0{
				path = fmt.Sprint(prefix,"/",strings.TrimPrefix(path,"/"))
			}

			service.server.Handle(path,&ViewServer{
				handler: view,
				modules: modules,
				name:    m.Name,
			})
			
		}

	}
	return service,nil
}

type ViewServer struct {
	handler ViewHandler
	modules interface{}
	name string
}

func (v *ViewServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	viewName, data, err := v.handler.Server(r)


	tp ,err:= template.Load(viewName)
	if err!=nil{
		fmt.Fprint(w,err)
		return
	}
	obj:=make(map[string]interface{})
	obj["data"]=data
	obj["error"] = err
	obj["modules"] = v.modules
	obj["current"] = v.name

	tp.Execute(w,obj)
}

type ModuleItem struct {
	Name string
	Icon string
	I18nName map[string]string
	Path string
}