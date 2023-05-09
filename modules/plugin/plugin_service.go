package plugin

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/plugin/plugin-model"
)

type IPluginService interface {
	GetList(ctx context.Context, namespaceId int) ([]*plugin_model.Plugin, error)
	Create(ctx context.Context, namespaceId, operator int, input *plugin_model.PluginInput) error
	Update(ctx context.Context, namespaceId, operator int, input *plugin_model.PluginInput) error
	Delete(ctx context.Context, namespaceId, operator int, name string) error
	GetByName(ctx context.Context, namespaceId int, name string) (*plugin_model.Plugin, error)
	GetBasicInfoList(ctx context.Context, namespaceId int) ([]*plugin_model.PluginBasic, error)  //只获取插件基本信息
	InsertBuilt(ctx context.Context, namespaceId int, plugins []*plugin_model.PluginBuilt) error //新增内置插件
	Sort(ctx context.Context, namespaceId, operator int, names []string) error
}
