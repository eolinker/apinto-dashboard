package main

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func main() {
	router := httprouter.New()
	h := func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		log.Println(params)
	}
	router.GET("/", h)
	router.GET("/profession/routers/", h)
	router.GET("/profession/routers/:driver", h)
	router.GET("/api/routers", h)
	router.GET("/api/routers/:name", h)

}
