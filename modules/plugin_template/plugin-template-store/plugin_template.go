package plugin_template_store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/plugin_template/plugin-template-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

var (
	_ IPluginTemplateStore = (*pluginTemplateStore)(nil)
)

type IPluginTemplateStore interface {
	store.IBaseStore[plugin_template_entry.PluginTemplate]
	GetListByNamespaceId(ctx context.Context, namespaceId int) ([]*plugin_template_entry.PluginTemplate, error)
	GetByUUID(ctx context.Context, uuid string) (*plugin_template_entry.PluginTemplate, error)
}

type pluginTemplateStore struct {
	*store.BaseStore[plugin_template_entry.PluginTemplate]
}

func newPluginTemplateStore(db store.IDB) IPluginTemplateStore {
	return &pluginTemplateStore{BaseStore: store.CreateStore[plugin_template_entry.PluginTemplate](db)}
}

func (p *pluginTemplateStore) GetListByNamespaceId(ctx context.Context, namespaceId int) ([]*plugin_template_entry.PluginTemplate, error) {
	list := make([]*plugin_template_entry.PluginTemplate, 0)
	if err := p.DB(ctx).Where("`namespace` = ?", namespaceId).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (p *pluginTemplateStore) GetByUUID(ctx context.Context, uuid string) (*plugin_template_entry.PluginTemplate, error) {
	val := new(plugin_template_entry.PluginTemplate)
	if err := p.DB(ctx).Where("`uuid` = ?", uuid).First(val).Error; err != nil {
		return nil, err
	}
	return val, nil
}
