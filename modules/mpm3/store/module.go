package store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/mpm3/entry"
	"github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/eosc/log"
)

type IModuleStore interface {
	GetEnable(ctx context.Context) []*entry.Module
	Save(ctx context.Context, plugin int, module []*entry.Module) error
}

func newPluginModuleStore(db store.IDB) IModuleStore {
	s := store.CreateStore[entry.Module](db)

	return &ModuleStore{base: s}
}

type ModuleStore struct {
	base store.IBaseStore[entry.Module]
}

func (m *ModuleStore) Save(ctx context.Context, plugin int, module []*entry.Module) error {
	return m.base.Transaction(ctx, func(ctx context.Context) error {
		_, err := m.base.DeleteWhere(ctx, map[string]interface{}{"plugin": plugin})
		if err != nil {
			return err
		}
		if len(module) == 0 {
			return nil
		}
		return m.base.Insert(ctx, module...)

	})
}

func (m *ModuleStore) GetEnable(ctx context.Context) []*entry.Module {
	ens := make([]*entry.Module, 0)
	err := m.base.DB(ctx).Table("pm3_module").Select("pm3_module.id,pm3_module.plugin,pm3_module.navigation,pm3_module.name,pm3_module.cname,pm3_module.router").
		Joins("join pm3_enable on pm3_enable.id = pm3_module.plugin and pm3_enable.is_enable=2").Scan(&ens).Error
	if err != nil {
		log.Error(err)
	}
	return ens
}
