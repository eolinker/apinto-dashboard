package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func main() {
	paths := []string{"/raft", "/rafts", "/raft/", "/raft/s"}
	router := httprouter.New()

	//router.HandlerFunc(http.MethodGet, "/raft/*more", func(w http.ResponseWriter, r *http.Request) {
	//	log.Println(r.URL)
	//})
	router.HandlerFunc(http.MethodGet, "/raft", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL)
		log.Println(r.RequestURI)

	})

	for _, p := range paths {
		h, params, b := router.Lookup(http.MethodGet, p)
		ok := "ok"
		if h == nil {
			ok = "none"
		}

		fmt.Printf("%s\t==>\t%s:%v:%v\n", p, ok, b, params)
	}

	http.ListenAndServe(":1000", router)
}
