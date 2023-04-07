package core

import (
	apinto_module "github.com/eolinker/apinto-module"
	"net/http"
)

type ICore interface {
	http.Handler
	ReloadModule(version string) error
}

type IProviders interface {
	apinto_module.IProviders
	Set(providers apinto_module.IProviders)
}
