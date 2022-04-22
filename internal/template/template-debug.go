//go:build dev
// +build dev

package template

import (
	"html/template"
)

func Load(name string) (*template.Template, error) {
	path := toPath(name)
	return read(path)
}
