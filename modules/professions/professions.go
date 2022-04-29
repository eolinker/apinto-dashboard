package professions

import (
	"fmt"
	apinto_dashboard "github.com/eolinker/apinto-dashboard"
	"github.com/eolinker/apinto-dashboard/internal/apinto"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type ListHeader struct {
	Title  map[apinto_dashboard.ZoneName][]string
	Fields []string
}
type Profession struct {
	*apinto_dashboard.ModuleViewFinder
	*httprouter.Router
	ModuleName     string
	ProfessionName string

	header *ListHeader
}

func (p *Profession) Paths() []string {
	paths := make([]string, 0, 2)
	paths = append(paths, fmt.Sprintf("/profession/%s/", p.ModuleName))
	paths = append(paths, fmt.Sprintf("/api/%s/", p.ModuleName))

	return paths
}

func (p *Profession) Lookup(r *http.Request) (view string, data interface{}, has bool) {
	name, has := p.ModuleViewFinder.Lookup(r)
	if has {
		switch name {
		case "profession_list":
			return name, p.header, true
		case "profession_edit":
			workerName := r.URL.Query().Get("name")
			return name, workerName, true
		}
		return name, nil, true
	}
	return "", nil, false
}

func (p *Profession) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	p.Router.ServeHTTP(w, req)
}

func NewProfession(name string, profession string, titles map[apinto_dashboard.ZoneName][]string, fields []string, viewFinder *apinto_dashboard.ModuleViewFinder) *Profession {

	if viewFinder == nil {
		views := map[string]string{
			"list":   "profession_list",
			"create": "profession_create",
			"edit":   "profession_edit",
		}
		viewFinder = apinto_dashboard.NewViewModuleEmpty(fmt.Sprint("/", name, "/"), views, "list")
	}
	p := &Profession{
		ModuleName:       name,
		ProfessionName:   profession,
		ModuleViewFinder: viewFinder,
		header: &ListHeader{
			Title:  titles,
			Fields: fields,
		},
	}
	p.createRouter()
	return p
}

func (p *Profession) createRouter() {
	r := httprouter.New()

	r.GET(fmt.Sprintf("/profession/%s/", p.ModuleName), func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		data, code, err := apinto.Client().Drivers(p.ProfessionName)
		if err != nil {
			apinto.WriteResult(w, 500, []byte(err.Error()))
			return
		}
		apinto.WriteResult(w, code, data)
	})
	r.GET(fmt.Sprintf("/profession/%s/:driver", p.ModuleName), func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		driver := params.ByName("driver")
		data, code, err := apinto.Client().Render(p.ProfessionName, driver)
		if err != nil {
			apinto.WriteResult(w, 500, []byte(err.Error()))
			return
		}
		apinto.WriteResult(w, code, data)
	})

	// List
	r.GET(fmt.Sprintf("/api/%s/", p.ModuleName), func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		data, code, err := apinto.Client().List(p.ProfessionName)
		if err != nil {
			apinto.WriteResult(w, 500, []byte(err.Error()))
			return
		}
		apinto.WriteResult(w, code, data)
	})

	// Get
	r.GET(fmt.Sprintf("/api/%s/:name", p.ModuleName), func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		name := params.ByName("name")
		data, code, err := apinto.Client().Get(p.ProfessionName, name)
		if err != nil {
			apinto.WriteResult(w, 500, []byte(err.Error()))
			return
		}
		apinto.WriteResult(w, code, data)
	})

	// Create
	r.POST(fmt.Sprintf("/api/%s/", p.ModuleName), func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		userName := apinto_dashboard.MustUsername(r)

		var logContent string
		logArgs := make([]*apinto_dashboard.Arg, 0, 6)
		logArgs = append(logArgs,
			&apinto_dashboard.Arg{Key: "user", Value: userName},
			&apinto_dashboard.Arg{Key: "profession", Value: p.ProfessionName},
			&apinto_dashboard.Arg{Key: "url", Value: r.URL.String()},
		)

		defer func() {
			apinto_dashboard.AddActivityLog(userName, "create", "", logContent, logArgs)
		}()

		rData, err := apinto.ReadBody(r.Body)
		if err != nil {
			logContent = fmt.Sprintf("编辑%s失败:Body读取失败", p.ProfessionName)
			logArgs = append(logArgs,
				&apinto_dashboard.Arg{Key: "error", Value: err.Error()},
				&apinto_dashboard.Arg{Key: "err_from", Value: "dashboard"},
			)

			apinto.WriteResult(w, 500, []byte(err.Error()))
			return
		}

		data, code, err := apinto.Client().Create(p.ProfessionName, rData)
		if err != nil {
			logContent = fmt.Sprintf("创建%s失败", p.ProfessionName)
			logArgs = append(logArgs,
				&apinto_dashboard.Arg{Key: "request_body", Value: string(rData)},
				&apinto_dashboard.Arg{Key: "error", Value: err.Error()},
				&apinto_dashboard.Arg{Key: "err_from", Value: "dashboard"},
			)

			apinto.WriteResult(w, 500, []byte(err.Error()))
		} else if code != 200 {
			logContent = fmt.Sprintf("创建%s失败", p.ProfessionName)
			logArgs = append(logArgs,
				&apinto_dashboard.Arg{Key: "request_body", Value: string(rData)},
				&apinto_dashboard.Arg{Key: "error", Value: string(data)},
				&apinto_dashboard.Arg{Key: "err_from", Value: "apinto"},
			)

			apinto.WriteResult(w, code, data)
		} else {
			logContent = fmt.Sprintf("创建%s成功", p.ProfessionName)
			logArgs = append(logArgs,
				&apinto_dashboard.Arg{Key: "request_body", Value: string(rData)},
			)

			apinto.WriteResult(w, code, data)
		}

	})

	// Update
	r.PUT(fmt.Sprintf("/api/%s/:name", p.ModuleName), func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		name := params.ByName("name")
		userName := apinto_dashboard.MustUsername(r)

		var logContent string
		logArgs := make([]*apinto_dashboard.Arg, 0, 6)
		logArgs = append(logArgs,
			&apinto_dashboard.Arg{Key: "user", Value: userName},
			&apinto_dashboard.Arg{Key: "profession", Value: p.ProfessionName},
			&apinto_dashboard.Arg{Key: "url", Value: r.URL.String()},
		)

		defer func() {
			apinto_dashboard.AddActivityLog(userName, "update", name, logContent, logArgs)
		}()

		rData, err := apinto.ReadBody(r.Body)
		if err != nil {
			logContent = fmt.Sprintf("编辑%s失败:Body读取失败", p.ProfessionName)
			logArgs = append(logArgs,
				&apinto_dashboard.Arg{Key: "error", Value: err.Error()},
				&apinto_dashboard.Arg{Key: "err_from", Value: "dashboard"},
			)

			apinto.WriteResult(w, 500, []byte(err.Error()))
			return
		}

		//更新
		data, code, err := apinto.Client().Update(p.ProfessionName, name, rData)
		if err != nil {
			logContent = fmt.Sprintf("编辑%s失败", p.ProfessionName)
			logArgs = append(logArgs,
				&apinto_dashboard.Arg{Key: "request_body", Value: string(rData)},
				&apinto_dashboard.Arg{Key: "error", Value: err.Error()},
				&apinto_dashboard.Arg{Key: "err_from", Value: "dashboard"},
			)

			apinto.WriteResult(w, 500, []byte(err.Error()))
		} else if code != 200 {
			logContent = fmt.Sprintf("编辑%s失败", p.ProfessionName)
			logArgs = append(logArgs,
				&apinto_dashboard.Arg{Key: "request_body", Value: string(rData)},
				&apinto_dashboard.Arg{Key: "error", Value: string(data)},
				&apinto_dashboard.Arg{Key: "err_from", Value: "apinto"},
			)

			apinto.WriteResult(w, code, data)
		} else {
			logContent = fmt.Sprintf("编辑%s成功", p.ProfessionName)
			logArgs = append(logArgs,
				&apinto_dashboard.Arg{Key: "request_body", Value: string(rData)},
			)

			apinto.WriteResult(w, code, data)
		}

	})

	// Delete
	r.DELETE(fmt.Sprintf("/api/%s/:name", p.ModuleName), func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		name := params.ByName("name")
		userName := apinto_dashboard.MustUsername(r)

		var logContent string
		logArgs := make([]*apinto_dashboard.Arg, 0, 5)
		logArgs = append(logArgs,
			&apinto_dashboard.Arg{Key: "user", Value: userName},
			&apinto_dashboard.Arg{Key: "profession", Value: p.ProfessionName},
			&apinto_dashboard.Arg{Key: "url", Value: r.URL.String()},
		)

		defer func() {
			apinto_dashboard.AddActivityLog(userName, "delete", name, logContent, logArgs)
		}()

		//删除
		data, code, err := apinto.Client().Delete(p.ProfessionName, name)
		if err != nil {
			logContent = fmt.Sprintf("删除%s失败", p.ProfessionName)
			logArgs = append(logArgs,
				&apinto_dashboard.Arg{Key: "error", Value: err.Error()},
				&apinto_dashboard.Arg{Key: "err_from", Value: "dashboard"},
			)

			apinto.WriteResult(w, 500, []byte(err.Error()))

		} else if code != 200 {
			logContent = fmt.Sprintf("删除%s失败", p.ProfessionName)
			logArgs = append(logArgs,
				&apinto_dashboard.Arg{Key: "error", Value: string(data)},
				&apinto_dashboard.Arg{Key: "err_from", Value: "apinto"},
			)

			apinto.WriteResult(w, code, data)
		} else {
			logContent = fmt.Sprintf("删除%s成功", p.ProfessionName)

			apinto.WriteResult(w, code, data)
		}

	})
	p.Router = r
}
