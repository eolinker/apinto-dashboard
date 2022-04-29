package template

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

var (
	appendView []string
	fileDir    string = "tpl"
)

func ResetAppendView(dir string, views []string) {
	fileDir = dir
	av := make([]string, 0, len(views))
	for _, v := range views {
		av = append(av, toPath(v))
	}
	appendView = av
}
func read(path string) (*template.Template, error) {
	ps := make([]string, 0, len(appendView)+1)
	ps = append(ps, appendView...)
	ps = append(ps, path)
	t, err := template.ParseFiles(ps...)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func toPath(name string) string {

	dir, _ := filepath.Abs(fileDir)
	if !strings.HasSuffix(name, ".html") {
		name = fmt.Sprint(name, ".html")
	}
	path := filepath.Join(dir, name)
	return path
}

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
