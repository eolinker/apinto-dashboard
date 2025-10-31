package service

import "github.com/eolinker/eosc/common/bean"

func init() {
	installService := newInstallService()
	bean.Injection(&installService)
	pluginService := NewPluginService()
	bean.Injection(&pluginService)
	moduleService := NewModuleService()
	bean.Injection(&moduleService)
	accessService := NewAccessService()
	bean.Injection(&accessService)

	frontendService := NewFrontendService()
	bean.Injection(&frontendService)

	resourcesService := NewResourcesService()
	bean.Injection(&resourcesService)
}
