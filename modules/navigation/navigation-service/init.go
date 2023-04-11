package navigation_service

import (
	module_plugin "github.com/eolinker/apinto-dashboard/modules/module-plugin"
	"github.com/eolinker/eosc/common/bean"
)

var (
	pluginService       module_plugin.IModulePlugin
	modulePluginService module_plugin.IModulePluginService
)

func init() {
	service := newNavigationService()
	bean.Injection(&service)
	bean.Autowired(&pluginService)
	bean.Autowired(&modulePluginService)
}
