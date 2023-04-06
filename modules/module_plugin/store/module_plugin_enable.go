package store

import (
	"github.com/eolinker/apinto-dashboard/modules/module_plugin/entry"
	"github.com/eolinker/apinto-dashboard/store"
)

var (
	_ IModulePluginEnableStore = (*modulePluginEnable)(nil)
)

type IModulePluginEnableStore interface {
	store.IBaseStore[entry.ModulePluginEnable]
}

type modulePluginEnable struct {
	*store.BaseStore[entry.ModulePluginEnable]
}

func newModulePluginEnableStore(db store.IDB) IModulePluginEnableStore {
	return &modulePluginEnable{BaseStore: store.CreateStore[entry.ModulePluginEnable](db)}
}
