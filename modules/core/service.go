package core

import "net/http"

type ICore interface {
	http.Handler
	ReloadModule(version string) error
}
