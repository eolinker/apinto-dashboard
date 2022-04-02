package apinto_dashboard

import (
	"bytes"
	"fmt"
	"github.com/eolinker/apinto-dashboard/internal/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

type IDashboardAccount interface {
}


type ViewHandlerFunc func(r *http.Request)(view string,data interface{},err error)

func (f ViewHandlerFunc) Lookup(r *http.Request) (view string, data interface{}, err error) {
	return f(r)
}

type Module struct {
	Path     string
	Handler IModule
	Name     string
	I18nName map[ZoneName]string `json:"i18n_name"`
}

type Config struct {
	DefaultZone ZoneName
	Modules []*Module
	UserDetailsService IUserDetailsService
	Statics map[string]string
	DefaultModule string
}
type DashboardService struct {
	defaultZone ZoneName
	userDetails IUserDetailsService
	serve http.ServeMux
	defaultModule string
}

func (d *DashboardService) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	d.serve.ServeHTTP(w,req)
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
			I18nName: m.I18nName,
			Path:     m.Path,
		})
	}
	mp := NewModuleItemPlan(modules)
	defaultModule:=config.DefaultModule
	if defaultModule == ""{
		defaultModule = modules[0].Name
	}
	views:=new(Views)
	views.mp = mp
	viewServe:=&http.ServeMux{}
	views.serve = viewServe
	apis := new(APIS)
	for _, m:=range config.Modules{

		path := fmt.Sprint("/", m.Name)
		viewServe.Handle(path,&ViewServer{
			handler: m.Handler,
			modules: mp,
			name:    m.Name,
		})
		service.serve.Handle(path,views)

		viewServe.Handle(fmt.Sprint(path,"/"),&ViewServer{
			handler: m.Handler,
			modules: mp,
			name:    m.Name,
		})
		service.serve.Handle(fmt.Sprint(path,"/"),views)
		apis.serve.Handle(fmt.Sprint("/api/",m.Name,"/"),m.Handler)
		service.serve.Handle(fmt.Sprint("/api/",m.Name,"/"),apis)
	}
	staticServe:=&http.ServeMux{}
	for path,dir:= range config.Statics{
		path = strings.TrimPrefix(path,"/")
		path = strings.TrimSuffix(path,"/")
		if len(path)>0{
			path = fmt.Sprint("/",path,"/")
			staticServe.Handle(path,http.StripPrefix(path,http.FileServer(http.Dir(dir))))

		}else{
			path = "/"
			staticServe.Handle(path,&Views{
				serve: http.StripPrefix(path,http.FileServer(http.Dir(dir))),
				mp: mp,
			})

		}
	}

	defaultModulePath := mp.moduleMap[defaultModule].Path
 	service.serve.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/"{
			http.Redirect(w,r,defaultModulePath,302)
			return
		}
		staticServe.ServeHTTP(w,r)
	})
	return service,nil
}

type APIS struct {
	serve http.ServeMux
}

func (A *APIS) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//todo 处理api接口的登陆
	A.serve.ServeHTTP(w,req)
}

type Views struct {
	serve http.Handler
	mp *ModuleItemPlan
}

func (v *Views) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//todo 处理view请求的登陆

	writer := NewTemplateWriter()

	v.serve.ServeHTTP(writer,req)

	if writer.statusCode == 200{
		writer.WriteTo(w)
		return
	}
	ext := filepath.Ext(req.URL.Path)
	if ext != ""{
		writer.WriteTo(w)
		return
	}
	tp, err :=template.Load("error")
	if err!= nil{
		log.Println("[ERR] load template<error>:",err)
		writer.WriteTo(w)
		return
	}
	tp.Execute(w,v.mp.CreateViewData("error",map[string]string{"statusCode":strconv.Itoa(writer.statusCode),"message":writer.buf.String()},nil))
}

type ViewServer struct {
	handler IModule
	modules *ModuleItemPlan
	name string
}

func (v *ViewServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	viewName, data, has := v.handler.Lookup(r)
	if !has{
		http.NotFound(w,r)
		return
	}
	tp ,err:= template.Load(viewName)
	if err!=nil{
		fmt.Fprint(w,err)
		return
	}

	tp.Execute(w,v.modules.CreateViewData(v.name,data,err))
}

type ModuleItem struct {
	Name string `json:"name"`
 	I18nName map[ZoneName]string `json:"i18n_name"`
	Path string `json:"path"`
}

type ModuleItemPlan struct {
	modules []*ModuleItem
	moduleMap map[string]*ModuleItem
}

func NewModuleItemPlan(modules []*ModuleItem) *ModuleItemPlan {
	mp:=make( map[string]*ModuleItem)
	for _,m:=range modules{
		mp[m.Name] = m
	}
	return &ModuleItemPlan{modules: modules,moduleMap: mp}
}
func (mp *ModuleItemPlan)CreateViewData(name string,data interface{},err error ) map[string]interface{} {
	obj:=make(map[string]interface{})
	obj["data"]=data
	obj["error"] = err
	obj["zone"] = ZhCn
	obj["modules"] = mp.modules
	obj["module"] = mp.moduleMap[name]
	obj["name"] = name
	return obj
}

type TemplateWriter struct {
	buf bytes.Buffer
	statusCode int
	header http.Header
}

func NewTemplateWriter() *TemplateWriter {
	return &TemplateWriter{
		statusCode: 200,
		header: make(http.Header),
	}
}
func (t *TemplateWriter) WriteTo(w http.ResponseWriter){


	for k:=range t.header{
		w.Header().Set(k,t.header.Get(k))
	}
	w.WriteHeader(t.statusCode)

	//w.Write(t.buf.Bytes())
	t.buf.WriteTo(w)
}
func (t *TemplateWriter) Header() http.Header {
	return t.header
}

func (t *TemplateWriter) Write(bytes []byte) (int, error) {
	return t.buf.Write(bytes)
}

func (t *TemplateWriter) WriteHeader(statusCode int) {
	t.statusCode = statusCode
}
