package local

import (
	apinto_module "github.com/eolinker/apinto-dashboard/module"
)

type rDriver struct {
}

func NewDriver() *rDriver {
	return &rDriver{}
}

func (d *rDriver) CreatePlugin(define interface{}) (apinto_module.Plugin, error) {
	return newRPlugin(define)
}
