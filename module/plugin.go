package apinto_module

type Plugin interface {
	CreateModule(name string, config interface{}) (Module, error)
	CheckConfig(name string, config interface{}) error
	GetPluginFrontend(moduleName string) string
	IsShowServer() bool
}
