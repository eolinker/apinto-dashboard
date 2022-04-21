package apinto_dashboard

import (
	"bytes"
	"fmt"
	"github.com/eolinker/apinto-dashboard/internal/template"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

type ViewHandlerFunc func(r *http.Request) (view string, data interface{}, err error)

func (f ViewHandlerFunc) Lookup(r *http.Request) (view string, data interface{}, err error) {
	return f(r)
}

type Module struct {
	Path     string
	Handler  IModule
	Name     string
	I18nName map[ZoneName]string `json:"i18n_name"`
}

type Config struct {
	DefaultZone        ZoneName
	Modules            []*Module
	UserDetailsService IUserDetailsService
	Statics            map[string]string
	DefaultModule      string
}

func Create(config *Config) (http.Handler, error) {
	if config.UserDetailsService == nil {
		return nil, ErrorUserDetailsServiceNeed
	}
	if config.DefaultZone == "" {
		config.DefaultZone = ZhCn
	}

	modules := make([]*ModuleItem, 0, len(config.Modules))
	for _, m := range config.Modules {
		modules = append(modules, &ModuleItem{
			Name:     m.Name,
			I18nName: m.I18nName,
			Path:     m.Path,
		})
	}
	mp := NewModuleItemPlan(modules)
	defaultModule := config.DefaultModule
	if defaultModule == "" {
		defaultModule = modules[0].Name
	}
	//views := new(Views)
	//views.mp = mp
	//views.serve = viewServe

	serve := &http.ServeMux{}
	for _, m := range config.Modules {

		path := fmt.Sprint("/", m.Name)
		viewH := &ViewServer{
			handler: m.Handler,
			modules: mp,
			name:    m.Name,
		}
		////viewServe.Handle(path, )
		//serve.Handle(path, viewH)
		//
		//viewServe.Handle(fmt.Sprint(path, "/"),viewH )
		//

		serve.Handle(path, viewH)
		serve.Handle(fmt.Sprint(path, "/"), viewH)
		serve.Handle(fmt.Sprint("/api/", m.Name, "/"), m.Handler)
		//serve.Handle(fmt.Sprint("/api/", m.Name, "/"), apis)
	}
	staticServe := &http.ServeMux{}

	for path, dir := range config.Statics {
		path = strings.TrimPrefix(path, "/")
		path = strings.TrimSuffix(path, "/")
		if len(path) > 0 {
			path = fmt.Sprint("/", path, "/")
			staticServe.Handle(path, http.StripPrefix(path, http.FileServer(http.Dir(dir))))
		} else {
			path = "/"
			//staticServe.Handle(path, &Views{
			//	serve: http.StripPrefix(path, http.FileServer(http.Dir(dir))),
			//	mp:    mp,
			//})
			staticServe.Handle(path, http.StripPrefix(path, http.FileServer(http.Dir(dir))))
		}
	}

	defaultModulePath := mp.moduleMap[defaultModule].Path
	serve.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, defaultModulePath, 302)
			return
		}
		staticServe.ServeHTTP(w, r)
	})
	return NewAccountHandler(config.UserDetailsService, &Views{
		serve: serve,
		mp:    mp,
	}, []string{"/css/", "/js/", "/umd/", "/fonts/"}), nil
}

type APIS struct {
	serve http.ServeMux
}

func (A *APIS) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//todo 处理api接口的登陆
	A.serve.ServeHTTP(w, req)
}

var viewExtHtml = map[string]int{
	".html": 1,
	".htm":  1,
	"":      1,
}

type Views struct {
	serve http.Handler
	mp    *ModuleItemPlan
}

func (v *Views) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	cache := NewTemplateWriter()

	v.serve.ServeHTTP(cache, req)
	codeHead := cache.statusCode / 100
	if codeHead != 4 && codeHead != 5 {
		cache.WriteTo(w)
		return
	}
	//accepts := req.Header.Get("Accept")
	//if strings.Contains(accepts, "text/html") || strings.Contains(accepts, "") {
	//	v.Error(w, cache)
	//	return
	//}
	if !strings.HasPrefix(req.URL.Path, "/api/") {
		ext := filepath.Ext(req.URL.Path)
		if viewExtHtml[ext] == 1 {
			v.Error(w, cache)
			return
		}
	}

	cache.WriteTo(w)
}
func (v *Views) Error(w http.ResponseWriter, cache *TemplateWriter) {
	template.Execute(w, "login", v.mp.CreateViewData("error", map[string]string{"statusCode": strconv.Itoa(cache.statusCode), "message": cache.buf.String()}, nil))

}

type ViewServer struct {
	handler IModule
	modules *ModuleItemPlan
	name    string
}

func (v *ViewServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	viewName, data, has := v.handler.Lookup(r)
	if !has {
		http.NotFound(w, r)
		return
	}
	template.Execute(w, viewName, v.modules.CreateViewData(v.name, data, nil))

}

type ModuleItem struct {
	Name     string              `json:"name"`
	I18nName map[ZoneName]string `json:"i18n_name"`
	Path     string              `json:"path"`
}

type ModuleItemPlan struct {
	modules   []*ModuleItem
	moduleMap map[string]*ModuleItem
}

func NewModuleItemPlan(modules []*ModuleItem) *ModuleItemPlan {
	mp := make(map[string]*ModuleItem)
	for _, m := range modules {
		mp[m.Name] = m
	}
	return &ModuleItemPlan{modules: modules, moduleMap: mp}
}
func (mp *ModuleItemPlan) CreateViewData(name string, data interface{}, err error) map[string]interface{} {
	obj := make(map[string]interface{})
	obj["data"] = data
	obj["error"] = err
	obj["zone"] = ZhCn
	obj["modules"] = mp.modules
	obj["module"] = mp.moduleMap[name]
	obj["name"] = name
	return obj
}

type TemplateWriter struct {
	buf        bytes.Buffer
	statusCode int
	header     http.Header
}

func NewTemplateWriter() *TemplateWriter {
	return &TemplateWriter{
		statusCode: 200,
		header:     make(http.Header),
	}
}
func (t *TemplateWriter) WriteTo(w http.ResponseWriter) {

	t.WriteHeaderTo(w)
	w.WriteHeader(t.statusCode)

	//w.Write(t.buf.Bytes())
	t.buf.WriteTo(w)
}
func (t *TemplateWriter) WriteHeaderTo(w http.ResponseWriter) {
	for k := range t.header {
		w.Header().Set(k, t.header.Get(k))
	}
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
