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
	r.GET("/api/:prfession/:name", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	})
	r.GET("/api/:prfession", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	})
	return &ProfessionRouter{}
}
