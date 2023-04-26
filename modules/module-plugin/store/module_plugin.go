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
	GetOtherGroupPlugins(ctx context.Context, groupIDs []string, searchName string) ([]*entry.PluginListItem, error)
	GetPluginList(ctx context.Context, groupID string, searchName string) ([]*entry.PluginListItem, error)
	GetInnerPluginList(ctx context.Context) ([]*entry.ModulePlugin, error)
	GetPluginInfo(ctx context.Context, uuid string) (*entry.ModulePlugin, error)
	GetEnabledPlugins(ctx context.Context) ([]*entry.EnablePlugin, error)
	GetNavigationModules(ctx context.Context) ([]*entry.EnabledModule, error)
}

type modulePluginStore struct {
	*store.BaseStore[entry.ModulePlugin]
}

func newModulePluginStore(db store.IDB) IModulePluginStore {
	return &modulePluginStore{BaseStore: store.CreateStore[entry.ModulePlugin](db)}
}

func (c *modulePluginStore) GetPluginList(ctx context.Context, groupID string, searchName string) ([]*entry.PluginListItem, error) {
	plugins := make([]*entry.PluginListItem, 0)

	db := c.DB(ctx).Table("module_plugin").Select("module_plugin.uuid, module_plugin.name, module_plugin.cname, module_plugin.resume, module_plugin.icon, module_plugin.type, module_plugin.group, module_plugin_enable.is_enable").
		Joins("right join module_plugin_enable on module_plugin.id = module_plugin_enable.id")
	if groupID != "" {
		db = db.Where("`group` = ?", groupID)
	}
	if searchName != "" {
		db = db.Where("`cname` like ? ", "%"+searchName+"%")
	}
	err := db.Order("create_time DESC").Scan(&plugins).Error
	return plugins, err
}

func (c *modulePluginStore) GetInnerPluginList(ctx context.Context) ([]*entry.ModulePlugin, error) {
	plugins := make([]*entry.ModulePlugin, 0)
	err := c.DB(ctx).Model(plugins).Where("`type` = 0 or `type` = 1 or `type` = 2").Find(&plugins).Error
	return plugins, err
}

func (c *modulePluginStore) GetOtherGroupPlugins(ctx context.Context, groupIDs []string, searchName string) ([]*entry.PluginListItem, error) {
	plugins := make([]*entry.PluginListItem, 0)

	db := c.DB(ctx).Table("module_plugin").Select("module_plugin.uuid, module_plugin.name, module_plugin.cname, module_plugin.resume, module_plugin.icon, module_plugin.type, module_plugin.group, module_plugin_enable.is_enable").
		Joins("right join module_plugin_enable on module_plugin.id = module_plugin_enable.id")
	if len(groupIDs) > 0 {
		db = db.Where("`group` not in (?)", groupIDs)
	}
	if searchName != "" {
		db = db.Where("`cname` like ? ", "%"+searchName+"%")
	}
	err := db.Order("create_time DESC").Scan(&plugins).Error
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

// GetNavigationModules 获取导航接口所需要的模块列表
func (c *modulePluginStore) GetNavigationModules(ctx context.Context) ([]*entry.EnabledModule, error) {
	modules := make([]*entry.EnabledModule, 0)
	err := c.DB(ctx).Table("module_plugin").Select("module_plugin_enable.name, module_plugin.cname, module_plugin.type, module_plugin_enable.navigation, module_plugin_enable.is_plugin_visible,module_plugin_enable.frontend").
		Joins("right join module_plugin_enable on module_plugin.id = module_plugin_enable.id").
		Where("module_plugin_enable.is_enable = 2").Scan(&modules).Error

	return modules, err
}
