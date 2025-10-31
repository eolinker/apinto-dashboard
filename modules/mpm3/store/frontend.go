package store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/mpm3/entry"
	"github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/eosc/log"
)

type IPluginFrontendStore interface {
	GetEnable(ctx context.Context) []*entry.Frontend
	Save(ctx context.Context, plugin int, module []*entry.Frontend) error
}

func newPluginFrontendStore(db store.IDB) IPluginFrontendStore {
	s := store.CreateStore[entry.Frontend](db)
	return &pluginFrontendStore{
		base: s,
	}
}

type pluginFrontendStore struct {
	base store.IBaseStore[entry.Frontend]
}

func (p *pluginFrontendStore) GetEnable(ctx context.Context) []*entry.Frontend {
	fs := make([]*entry.Frontend, 0)
	err := p.base.DB(ctx).Table("pm3_frontend").Select("pm3_frontend.id,pm3_frontend.plugin,pm3_frontend.content").
		Joins("join pm3_enable on pm3_enable.id = pm3_frontend.plugin and pm3_enable.is_enable=2").Scan(&fs).Error
	if err != nil {
		log.Error(err)
	}
	return fs
}

func (p *pluginFrontendStore) Save(ctx context.Context, plugin int, fs []*entry.Frontend) error {
	_, err := p.base.DeleteWhere(ctx, map[string]interface{}{"plugin": plugin})
	if err != nil {
		return err
	}
	if len(fs) == 0 {
		return nil
	}
	return p.base.Insert(ctx, fs...)
}
