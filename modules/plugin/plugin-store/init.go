package plugin_store

import (
	"github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/eosc/common/bean"
)

func init() {
	store.RegisterStore(func(db store.IDB) {
		plugin := newPluginStore(db)
		pluginHistory := newPluginHistoryStore(db)

		clusterPlugin := newClusterPluginStore(db)
		clusterPluginHistory := newClusterPluginHistoryStore(db)
		clusterPluginPublishHistory := newClusterPluginPublishHistoryStore(db)
		clusterPluginPublishVersion := newClusterPluginPublishVersionStore(db)
		clusterPluginRuntime := newClusterPluginRuntimeStore(db)

		bean.Injection(&plugin)
		bean.Injection(&pluginHistory)
		bean.Injection(&clusterPlugin)
		bean.Injection(&clusterPluginHistory)
		bean.Injection(&clusterPluginPublishHistory)
		bean.Injection(&clusterPluginRuntime)
		bean.Injection(&clusterPluginPublishVersion)
	})
}
