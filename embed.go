package apinto_dashboard

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"
)

//go:embed builds/static
var staticDir embed.FS

func getStaticFiles() http.FileSystem {
	files, err := fs.Sub(staticDir, "builds/static")
	if err != nil {
		panic(err)
	}
	return http.FS(files)
}

//go:embed builds/tpl
var tplDir embed.FS

var templateFileFunc func(path string) (*template.Template, error)

func getTemplateFile(path string) (*template.Template, error) {
	bytes, err := tplDir.ReadFile(path)
	if err != nil {
		return nil, err
	}
	t := template.New(path)

	t2, err := t.Parse(string(bytes))
	if err != nil {
		return nil, err
	}

	return t2, nil
}
