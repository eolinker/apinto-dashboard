package service

import (
	"github.com/eolinker/eosc/common/bean"
)

func init() {
	iModulePluginService := newModulePluginService()
	iModulePlugin := newModulePlugin()
	bean.Injection(&iModulePluginService)
	bean.Injection(&iModulePlugin)
	iInstalledCache := newIInstalledCache()
	iNavigationModulesCache := newNavigationModulesCache()
	bean.Injection(&iInstalledCache)
	bean.Injection(&iNavigationModulesCache)

}
