//go:build !dev
// +build !dev

package template

import (
	"html/template"
	"sync"
)

var (
	cache = make(map[string]*template.Template)
	lock  sync.RWMutex
)

func Load(name string) (*template.Template, error) {
	path := toPath(name)
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
	t, err := read(path)
	if err != nil {
		return nil, err
	}
	cache[path] = t
	return t, nil

}
