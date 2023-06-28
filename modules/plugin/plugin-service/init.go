package plugin_service

import (
	"github.com/eolinker/eosc/common/bean"
)

func init() {
	pluginServiceInfo := newPluginService()
	clusterPluginServiceInfo := newClusterPluginService()
	bean.Injection(&pluginServiceInfo)
	bean.Injection(&clusterPluginServiceInfo)
	iExtenderCache := newIExtenderCache()
	bean.Injection(&iExtenderCache)

}
