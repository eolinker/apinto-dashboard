package apinto

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

var (
	professionRouterClientHandler *ProfessionRouter = NewProfessionRouter()
)

func Handler() http.Handler {
	return professionRouterClientHandler
}

type ProfessionRouter struct {
	*httprouter.Router
}

func NewProfessionRouter() *ProfessionRouter {
	r := httprouter.New()

	// Get
	r.GET("/api/:profession/:name", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		profession := params.ByName("profession")
		name := params.ByName("name")
		resp, err := Client().Get(profession, name)
		if err != nil {
			writeResult(w, 500, []byte(err.Error()))
			return
		}
		writeResult(w, resp.code, resp.data)
	})

	// List
	r.GET("/api/:profession", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		profession := params.ByName("profession")
		resp, err := Client().List(profession)
		if err != nil {
			writeResult(w, 500, []byte(err.Error()))
			return
		}
		writeResult(w, resp.code, resp.data)
	})

	// Create
	r.POST("/api/:profession", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		profession := params.ByName("profession")
		data, err := readBody(r.Body)
		if err != nil {
			writeResult(w, 500, []byte(err.Error()))
			return
		}
		resp, err := Client().Create(profession, data)
		if err != nil {
			writeResult(w, 500, []byte(err.Error()))
			return
		}
		writeResult(w, resp.code, resp.data)
	})

	// Update
	r.PUT("/api/:profession/:name", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		profession := params.ByName("profession")
		name := params.ByName("name")
		data, err := readBody(r.Body)
		if err != nil {
			writeResult(w, 500, []byte(err.Error()))
			return
		}
		resp, err := Client().Update(profession, name, data)
		if err != nil {
			writeResult(w, 500, []byte(err.Error()))
			return
		}
		writeResult(w, resp.code, resp.data)
	})

	// Delete
	r.DELETE("/api/:profession/:name", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		profession := params.ByName("profession")
		name := params.ByName("name")
		resp, err := Client().Delete(profession, name)
		if err != nil {
			writeResult(w, 500, []byte(err.Error()))
			return
		}
		writeResult(w, resp.code, resp.data)
	})

	// Patch
	r.PATCH("/api/:profession/:name", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		profession := params.ByName("profession")
		name := params.ByName("name")
		data, err := readBody(r.Body)
		if err != nil {
			writeResult(w, 500, []byte(err.Error()))
			return
		}
		resp, err := Client().Patch(profession, name, data)
		if err != nil {
			writeResult(w, 500, []byte(err.Error()))
			return
		}
		writeResult(w, resp.code, resp.data)
	})

	// Enable
	r.PATCH("/api/:profession/:name/status", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		profession := params.ByName("profession")
		name := params.ByName("name")
		if profession != "router" {
			writeResult(w, 200, []byte("just router can change status"))
			return
		}
		resp, err := Client().Enable(profession, name)
		if err != nil {
			writeResult(w, 500, []byte(err.Error()))
			return
		}
		writeResult(w, resp.code, resp.data)
	})
	return &ProfessionRouter{}
}
