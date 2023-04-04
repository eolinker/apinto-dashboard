package store

import (
	"github.com/eolinker/apinto-dashboard/modules/module_plugin/entry"
	"github.com/eolinker/apinto-dashboard/store"
)

var (
	_ IModulePluginPackageStore = (*modulePluginPackage)(nil)
)

type IModulePluginPackageStore interface {
	store.IBaseStore[entry.ModulePluginPackage]
}

type modulePluginPackage struct {
	*store.BaseStore[entry.ModulePluginPackage]
}

func newModulePluginPackageStore(db store.IDB) IModulePluginPackageStore {
	return &modulePluginPackage{BaseStore: store.CreateStore[entry.ModulePluginPackage](db)}
}
