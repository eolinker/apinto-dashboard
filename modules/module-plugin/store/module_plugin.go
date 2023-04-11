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
	GetEnabledPlugins(ctx context.Context) ([]*entry.EnablePlugin, error)
	GetEnabledModules(ctx context.Context, navigationIds []int) ([]*entry.EnabledModule, error)
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

func (c *modulePluginStore) GetEnabledPlugins(ctx context.Context) ([]*entry.EnablePlugin, error) {
	plugins := make([]*entry.EnablePlugin, 0)
	err := c.DB(ctx).Table("module_plugin").Select("module_plugin.uuid, module_plugin_enable.name, module_plugin.driver, module_plugin_enable.config, module_plugin.details").
		Joins("right join module_plugin_enable on module_plugin.id = module_plugin_enable.id").
		Where("module_plugin_enable.is_enable = 2").
		Scan(&plugins).Error

	return plugins, err
}

// GetEnabledModules 通过导航id列表获取启用的模块信息
func (c *modulePluginStore) GetEnabledModules(ctx context.Context, navigationIds []int) ([]*entry.EnabledModule, error) {
	modules := make([]*entry.EnabledModule, 0)
	db := c.DB(ctx).Table("module_plugin").Select("module_plugin_enable.name, module_plugin.cname, module_plugin.type, module_plugin_enable.navigation, module_plugin.front").
		Joins("right join module_plugin_enable on module_plugin.id = module_plugin_enable.id")
	if len(navigationIds) > 0 {
		db = db.Where("module_plugin_enable.navigation in (?)", navigationIds)
	}
	err := db.Where("module_plugin_enable.is_enable = 2").Scan(&modules).Error

	return modules, err
}
