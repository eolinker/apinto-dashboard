//go:build !dev
// +build !dev

package template

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"strings"
	"sync"
)

var (
	cache   = make(map[string]*template.Template)
	lock    sync.RWMutex
	tplFile fs.FS
)

//go:embed tpl
var tplDir embed.FS

func init() {
	sub, err := fs.Sub(tplDir, "tpl")
	if err != nil {
		panic(err)
	}
	tplFile = sub
}

func ResetAppendView(dir string, views []string) {

	av := make([]string, 0, len(views))
	for _, v := range views {
		av = append(av, toEmbedPath(v))
	}
	appendView = av
}
func toEmbedPath(name string) string {

	if !strings.HasSuffix(name, ".html") {
		name = fmt.Sprint(name, ".html")
	}

	return name
}
func Load(name string) (*template.Template, error) {
	path := toEmbedPath(name)
	if tm, has := readCache(path); has {
		return tm, nil
	}
	return load(path)
}

func readCache(path string) (*template.Template, bool) {
	lock.RLock()
	tm, has := cache[path]
	lock.RUnlock()
	return tm, has
}
func load(path string) (*template.Template, error) {
	lock.Lock()
	defer lock.Unlock()
	tm, has := cache[path]
	if has {
		return tm, nil
	}
	t, err := readEmbed(path)
	if err != nil {
		return nil, err
	}
	cache[path] = t
	return t, nil

}

func readEmbed(path string) (*template.Template, error) {
	ps := make([]string, 0, len(appendView)+1)
	ps = append(ps, appendView...)
	ps = append(ps, path)

	t, err := template.ParseFS(tplFile, ps...)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return t, nil
}
