package plugin_store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/plugin/plugin-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

var (
	_ IPluginStore = (*pluginStore)(nil)
)

type IPluginStore interface {
	store.IBaseStore[plugin_entry.Plugin]
	GetListByNamespaceId(ctx context.Context, namespaceId int) ([]*plugin_entry.Plugin, error)
	GetByName(ctx context.Context, namespaceId int, name string) (*plugin_entry.Plugin, error)
	MaxSort(ctx context.Context, namespaceId int) (int, error)
}

type pluginStore struct {
	*store.BaseStore[plugin_entry.Plugin]
}

func newPluginStore(db store.IDB) IPluginStore {
	return &pluginStore{BaseStore: store.CreateStore[plugin_entry.Plugin](db)}
}

func (p *pluginStore) GetListByNamespaceId(ctx context.Context, namespaceId int) ([]*plugin_entry.Plugin, error) {
	list := make([]*plugin_entry.Plugin, 0)
	if err := p.DB(ctx).Where("`namespace` = ?", namespaceId).Order("`sort` asc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (p *pluginStore) GetByName(ctx context.Context, namespaceId int, name string) (*plugin_entry.Plugin, error) {
	val := new(plugin_entry.Plugin)
	if err := p.DB(ctx).Where("`namespace` = ? and `name` = ?", namespaceId, name).First(val).Error; err != nil {
		return nil, err
	}
	return val, nil
}

func (p *pluginStore) MaxSort(ctx context.Context, namespaceId int) (int, error) {
	sort := 0

	db := p.DB(ctx).Table("`plugin`").Select("IFNULL(MAX(`sort`),0) AS `sort`") //IFNULL MAX 为了处理 N/A默认值的问题

	if err := db.Where("`namespace` = ?", namespaceId).Order("`sort` desc").Limit(1).Row().Scan(&sort); err != nil {
		return 0, err
	}

	return sort, nil
}
