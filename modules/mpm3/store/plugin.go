package store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/mpm3/entry"
	"github.com/eolinker/apinto-dashboard/store"
)

var (
	_ IPluginStore = (*modulePluginStore)(nil)
)

type IPluginStore interface {
	store.IBaseStore[entry.Plugin]
	SearchFromOther(ctx context.Context, groupIDs []string, searchName string) ([]*entry.PluginListItem, error)
	Search(ctx context.Context, groupID string, searchName string) ([]*entry.PluginListItem, error)
	GetInnerPluginList(ctx context.Context) ([]*entry.Plugin, error)
	GetPluginInfo(ctx context.Context, uuid string) (*entry.Plugin, error)
	GetEnabledPlugins(ctx context.Context) ([]*entry.EnablePlugin, error)
	DeleteNotIn(ctx context.Context, inert bool, uuids ...string) ([]int, error)
}

type modulePluginStore struct {
	*store.BaseStore[entry.Plugin]
}

func (c *modulePluginStore) DeleteNotIn(ctx context.Context, inert bool, uuids ...string) ([]int, error) {
	deleteIds := make([]int, 0)
	err := c.DB(ctx).Table("pm3_plugins").Select("pm3_plugins.id").Joins("left join pm3_enable on pm3_plugins.id = pm3_enable.id").Where("pm3_plugins.is_inner =? and  pm3_plugins.uuid not in(?)", inert, uuids).Scan(&deleteIds).Error
	if err != nil {
		return nil, err
	}
	if len(deleteIds) == 0 {
		return nil, nil
	}
	_, err = c.BaseStore.Delete(ctx, deleteIds...)
	if err != nil {
		return nil, err
	}

	return deleteIds, err
}

func newModulePluginStore(db store.IDB) IPluginStore {
	ctx := context.Background()
	ms := &modulePluginStore{BaseStore: store.CreateStore[entry.Plugin](db)}
	migrate := db.DB(ctx).Migrator()
	if migrate.HasColumn(&entry.Plugin{}, "type") {
		db.DB(ctx).Model(&entry.Plugin{}).Where("type = 3").Updates(map[string]interface{}{
			"is_can_disable":        1,
			"is_can_uninstall":      1,
			"visible_in_navigation": 1,
			"visible_in_market":     1,
		})

		db.DB(ctx).Model(&entry.Plugin{}).Where("type < 3").Updates(map[string]interface{}{
			"is_inner": 1,
		})

		err := migrate.DropColumn(&entry.Plugin{}, "type")
		if err != nil {
			panic(err)
		}
	}
	return ms
}

func (c *modulePluginStore) Search(ctx context.Context, groupID string, searchName string) ([]*entry.PluginListItem, error) {
	plugins := make([]*entry.PluginListItem, 0)

	db := c.DB(ctx).Table("pm3_plugins").Select("pm3_plugins.uuid, pm3_plugins.name, pm3_plugins.cname, pm3_plugins.resume, " +
		"pm3_plugins.icon, pm3_plugins.is_inner,  pm3_plugins.group, pm3_enable.is_enable").
		Joins("left join pm3_enable on pm3_plugins.id = pm3_enable.id")
	if groupID != "" {
		db = db.Where("`group` = ?", groupID)
	}
	if searchName != "" {
		db = db.Where("`cname` like ? ", "%"+searchName+"%")
	}
	err := db.Order("create_time DESC").Scan(&plugins).Error
	return plugins, err
}

func (c *modulePluginStore) GetInnerPluginList(ctx context.Context) ([]*entry.Plugin, error) {
	plugins := make([]*entry.Plugin, 0)
	err := c.DB(ctx).Model(plugins).Where("`is_inner` = 1").Find(&plugins).Error
	return plugins, err
}

func (c *modulePluginStore) SearchFromOther(ctx context.Context, groupIDs []string, searchName string) ([]*entry.PluginListItem, error) {
	plugins := make([]*entry.PluginListItem, 0)

	db := c.DB(ctx).Table("pm3_plugins").Select("pm3_plugins.uuid, pm3_plugins.name, pm3_plugins.cname, pm3_plugins.resume, pm3_plugins.icon, " +
		"pm3_plugins.is_inner,  pm3_plugins.group, pm3_enable.is_enable").
		Joins("right join pm3_enable on pm3_plugins.id = pm3_enable.id")
	if len(groupIDs) > 0 {
		db = db.Where("`group` not in (?)", groupIDs)
	}
	if searchName != "" {
		db = db.Where("`cname` like ? ", "%"+searchName+"%")
	}
	err := db.Order("create_time DESC").Scan(&plugins).Error
	return plugins, err
}

func (c *modulePluginStore) GetPluginInfo(ctx context.Context, uuid string) (*entry.Plugin, error) {
	return c.FirstQuery(ctx, "`uuid` = ?", []interface{}{uuid}, "")
}

func (c *modulePluginStore) GetEnabledPlugins(ctx context.Context) ([]*entry.EnablePlugin, error) {
	plugins := make([]*entry.EnablePlugin, 0)
	err := c.DB(ctx).Table("pm3_plugins").Select("pm3_plugins.uuid,  pm3_plugins.driver, pm3_enable.config, pm3_plugins.details").
		Joins("right join pm3_enable on pm3_plugins.id = pm3_enable.id").
		Where("pm3_enable.is_enable = 2").
		Scan(&plugins).Error

	return plugins, err
}
