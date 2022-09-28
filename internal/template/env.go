package template

import (
	"fmt"
	"log"
	"net/http"
)

var (
	appendView []string
)

func Execute(w http.ResponseWriter, tplName string, data interface{}) {
	tp, err := Load(tplName)
	if err != nil {
		log.Println("[ERR] load template<error>:", err)
		w.WriteHeader(504)
		fmt.Fprintln(w, "[ERR] load template<error>:", err)
		return
	}
	tp.Execute(w, data)
	return
}
