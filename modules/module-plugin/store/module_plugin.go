package store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/entry"
	"github.com/eolinker/apinto-dashboard/store"
)

var (
	_ IModulePluginStore = (*modulePluginStore)(nil)
)

type IModulePluginStore interface {
	store.IBaseStore[entry.ModulePlugin]
	GetPluginList(ctx context.Context, groupUUID string, pluginName string) ([]*entry.ModulePlugin, error)
}

type modulePluginStore struct {
	*store.BaseStore[entry.ModulePlugin]
}

func newModulePluginStore(db store.IDB) IModulePluginStore {
	return &modulePluginStore{BaseStore: store.CreateStore[entry.ModulePlugin](db)}
}

func (c *modulePluginStore) GetPluginList(ctx context.Context, groupUUID string, pluginName string) ([]*entry.ModulePlugin, error) {
	plugins := make([]*entry.ModulePlugin, 0)
	db := c.DB(ctx).Model(plugins)
	if groupUUID != "" {
		db = db.Where("`group` = ?", groupUUID)
	}
	if pluginName != "" {
		db = db.Where("`cname` like ? ", "%"+pluginName+"%")
	}
	err := db.Order("create_time DESC").Find(&plugins).Error
	return plugins, err
}
