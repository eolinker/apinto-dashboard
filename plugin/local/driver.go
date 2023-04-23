package local

import (
	apinto_module "github.com/eolinker/apinto-module"
)

type tDriver struct {
}

func NewDriver() *tDriver {
	return &tDriver{}
}

func (d *tDriver) CreatePlugin(define interface{}) (apinto_module.Plugin, error) {
	return newTPlugin(define)
}
