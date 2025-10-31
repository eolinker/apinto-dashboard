package store

import (
	"github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/eosc/common/bean"
)

func init() {
	store.RegisterStore(func(db store.IDB) {
		modulePlugin := newModulePluginStore(db)
		bean.Injection(&modulePlugin)

		moduleStory := newPluginModuleStore(db)
		bean.Injection(&moduleStory)

		pluginEnable := newModulePluginEnableStore(db)
		bean.Injection(&pluginEnable)

		modulePluginResources := newPluginResourcesStore(db, modulePlugin)
		bean.Injection(&modulePluginResources)

		frontendStore := newPluginFrontendStore(db)
		bean.Injection(&frontendStore)
		accessStore := newPluginAccessStore(db)
		bean.Injection(&accessStore)
	})
}
