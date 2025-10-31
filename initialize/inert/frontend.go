package inert

import "github.com/eolinker/apinto-dashboard/pm3"

var (
	inertPluginFrontendConfig []pm3.PFrontend
)

func AddInertFrontendConfig(fs ...pm3.PFrontend) {
	inertPluginFrontendConfig = append(inertPluginFrontendConfig, fs...)
}

func GetFrontends() []pm3.PFrontend {
	return inertPluginFrontendConfig
}
