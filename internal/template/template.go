package template

import (
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
	"sync"
)

var (
	indexView = "index"
	fileDir string = "tpl"
	cache = make(map[string]*template.Template)
	lock sync.RWMutex
)

func Load(name string) (*template.Template,error) {

	path:=toPath(name)
	if tm,has:=readCache(path);has{
		return tm,nil
	}
	return load(indexView,path)
}

func readCache(path string)(*template.Template,bool)  {
	lock.RLock()
	tm,has:=cache[path]
	lock.RUnlock()
	return tm,has
}
func load(index,path string) (*template.Template,error) {
	lock.Lock()
	defer lock.Unlock()
	tm,has:=cache[path]
	if has{
		return tm,nil
	}
	t, err := template.ParseFiles(toPath(index),toPath("error"),path)
	if err!= nil{
		return nil,err
	}

	cache[path]=t
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