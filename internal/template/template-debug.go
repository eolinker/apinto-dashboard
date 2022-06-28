//go:build dev
// +build dev

package template

import (
	"html/template"
)

var (
	fileDir string = "tpl"
)

func ResetAppendView(dir string, views []string) {
	fileDir = dir
	av := make([]string, 0, len(views))
	for _, v := range views {
		av = append(av, toPath(v))
	}
	appendView = av
}

func Load(name string) (*template.Template, error) {
	path := toPath(name)
	return read(path)
}

func toPath(name string) string {

	dir, _ := filepath.Abs(fileDir)
	if !strings.HasSuffix(name, ".html") {
		name = fmt.Sprint(name, ".html")
	}
	path := filepath.Join(dir, name)
	return path
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
