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
	GetPluginList(ctx context.Context, groupID int, searchName string) ([]*entry.ModulePlugin, error)
	GetPluginInfo(ctx context.Context, uuid string) (*entry.ModulePlugin, error)
}

type modulePluginStore struct {
	*store.BaseStore[entry.ModulePlugin]
}

func newModulePluginStore(db store.IDB) IModulePluginStore {
	return &modulePluginStore{BaseStore: store.CreateStore[entry.ModulePlugin](db)}
}

func (c *modulePluginStore) GetPluginList(ctx context.Context, groupID int, searchName string) ([]*entry.ModulePlugin, error) {
	plugins := make([]*entry.ModulePlugin, 0)
	db := c.DB(ctx).Model(plugins)
	if groupID > 0 {
		db = db.Where("`group` = ?", groupID)
	}
	if searchName != "" {
		db = db.Where("`cname` like ? ", "%"+searchName+"%")
	}
	err := db.Order("create_time DESC").Find(&plugins).Error
	return plugins, err
}

func (c *modulePluginStore) GetPluginInfo(ctx context.Context, uuid string) (*entry.ModulePlugin, error) {
	return c.FirstQuery(ctx, "`uuid` = ?", []interface{}{uuid}, "")
}
