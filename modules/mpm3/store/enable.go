package store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/mpm3/entry"
	"github.com/eolinker/apinto-dashboard/store"
)

var (
	_ EnableStore = (*modulePluginEnable)(nil)
)

type EnableStore interface {
	//store.IBaseStore[entry.PluginEnable]
	Get(ctx context.Context, id int) (*entry.PluginEnable, error)
	Delete(ctx context.Context, id ...int) (int, error)
	Save(ctx context.Context, enable *entry.PluginEnable) error
}

func (c *modulePluginEnable) DeleteNotIn(ctx context.Context, uuids ...string) (int, error) {
	r := c.base.DB(ctx).Table("pm3_enable").Where(" is_enable != 2  and `id` not in(select id from pm3_plugin where uuid not in(?)", uuids).Delete("pm3_enable")
	return int(r.RowsAffected), r.Error
}

type modulePluginEnable struct {
	base *store.BaseStore[entry.PluginEnable]
}

func (c *modulePluginEnable) Get(ctx context.Context, id int) (*entry.PluginEnable, error) {
	return c.base.Get(ctx, id)
}

func (c *modulePluginEnable) Delete(ctx context.Context, id ...int) (int, error) {
	return c.base.Delete(ctx, id...)
}

func (c *modulePluginEnable) Save(ctx context.Context, enable *entry.PluginEnable) error {
	return c.base.Save(ctx, enable)
}

func newModulePluginEnableStore(db store.IDB) EnableStore {
	migrate := db.DB(context.Background()).Migrator()
	if migrate.HasColumn(&entry.PluginEnable{}, "is_can_disable") {
		err := migrate.DropColumn(&entry.PluginEnable{}, "is_can_disable")
		if err != nil {
			panic(err)
		}
	}
	if migrate.HasColumn(&entry.PluginEnable{}, "is_can_uninstall") {
		err := migrate.DropColumn(&entry.PluginEnable{}, "is_can_uninstall")
		if err != nil {
			panic(err)
		}
	}
	if migrate.HasColumn(&entry.PluginEnable{}, "is_show_server") {
		err := migrate.DropColumn(&entry.PluginEnable{}, "is_show_server")
		if err != nil {
			panic(err)
		}
	}
	if migrate.HasColumn(&entry.PluginEnable{}, "is_plugin_visible") {
		err := migrate.DropColumn(&entry.PluginEnable{}, "is_plugin_visible")
		if err != nil {
			panic(err)
		}
	}
	return &modulePluginEnable{base: store.CreateStore[entry.PluginEnable](db)}
}
