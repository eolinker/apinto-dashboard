package store

import (
	"github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/eosc/common/bean"
)

func init() {
	store.RegisterStore(func(db store.IDB) {
		modulePlugin := newModulePluginStore(db)
		pluginEnable := newModulePluginEnableStore(db)
		pluginPackage := newModulePluginPackageStore(db)

		bean.Injection(&modulePlugin)
		bean.Injection(&pluginEnable)
		bean.Injection(&pluginPackage)
	})
}
