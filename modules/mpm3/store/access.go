package store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/mpm3/entry"
	"github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/eosc/log"
)

type IPluginAccessStore interface {
	GetEnable(ctx context.Context) []*entry.Access
	Save(ctx context.Context, plugin int, module []*entry.Access) error
}

func newPluginAccessStore(db store.IDB) IPluginAccessStore {
	s := store.CreateStore[entry.Access](db)
	return &pluginAccessStore{
		base: s,
	}
}

type pluginAccessStore struct {
	base store.IBaseStore[entry.Access]
}

func (p *pluginAccessStore) GetEnable(ctx context.Context) []*entry.Access {
	accessList := make([]*entry.Access, 0)
	err := p.base.DB(ctx).Table("pm3_access").Select("pm3_access.id,pm3_access.name,pm3_access.cname,pm3_access.plugin,pm3_access.module,pm3_access.depend").
		Joins("inner join pm3_enable on pm3_enable.id = pm3_access.plugin and pm3_enable.is_enable=2").Scan(&accessList).Error
	if err != nil {
		log.Error(err)
	}
	return accessList
}

func (p *pluginAccessStore) Save(ctx context.Context, plugin int, module []*entry.Access) error {
	_, err := p.base.DeleteWhere(ctx, map[string]interface{}{"plugin": plugin})
	if err != nil {
		return err
	}
	if len(module) == 0 {
		return nil
	}
	return p.base.Insert(ctx, module...)
}
