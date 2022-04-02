// +build dev

package template

import (
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
)

var (
	indexView = "index"
	fileDir   = "tpl"
)

func Load(name string) (*template.Template,error) {
	path:=toPath(name)

	return load(indexView,path)
}

func load(index,path string) (*template.Template,error) {

	t, err := template.ParseFiles(toPath(index),path)
	if err!= nil{
		return nil,err
	}

	return t,nil

}

func toPath(name string)string  {

	dir,_:=filepath.Abs(fileDir)
	if !strings.HasSuffix(name,".html"){
		name = fmt.Sprint(name,".html")
	}
	path:= filepath.Join(dir,name)
	return path
}