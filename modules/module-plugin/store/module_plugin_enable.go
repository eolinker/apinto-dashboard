package store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/entry"
	"github.com/eolinker/apinto-dashboard/store"
)

var (
	_ IModulePluginEnableStore = (*modulePluginEnable)(nil)
)

type IModulePluginEnableStore interface {
	store.IBaseStore[entry.ModulePluginEnable]
	GetListByNavigation(ctx context.Context, navigationID int) ([]*entry.ModulePluginEnable, error)
	GetEnabledPluginByName(ctx context.Context, moduleName string) (*entry.ModulePluginEnable, error)
}

type modulePluginEnable struct {
	*store.BaseStore[entry.ModulePluginEnable]
}

func newModulePluginEnableStore(db store.IDB) IModulePluginEnableStore {
	return &modulePluginEnable{BaseStore: store.CreateStore[entry.ModulePluginEnable](db)}
}

func (m *modulePluginEnable) GetListByNavigation(ctx context.Context, navigationID int) ([]*entry.ModulePluginEnable, error) {
	return m.List(ctx, map[string]interface{}{
		"navigation": navigationID,
	})
}

func (m *modulePluginEnable) GetEnabledPluginByName(ctx context.Context, moduleName string) (*entry.ModulePluginEnable, error) {
	return m.FirstQuery(ctx, "name = ? and is_enable = 2", []interface{}{moduleName}, "")
}
