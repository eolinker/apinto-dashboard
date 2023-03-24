package plugin_template

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/base/frontend-model"
	"github.com/eolinker/apinto-dashboard/modules/plugin_template/plugin-template-model"
)

type IPluginTemplateService interface {
	GetList(ctx context.Context, namespaceId int) ([]*plugin_template_model.PluginTemplate, error)
	GetUsableList(ctx context.Context, namespaceId int) ([]*plugin_template_model.PluginTemplate, error)
	Create(ctx context.Context, namespaceId, operator int, input *plugin_template_model.PluginTemplateDetail) error
	Update(ctx context.Context, namespaceId, operator int, input *plugin_template_model.PluginTemplateDetail) error
	Delete(ctx context.Context, namespaceId, operator int, uuid string) error
	GetByUUID(ctx context.Context, namespaceId int, uuid string) (*plugin_template_model.PluginTemplateDetail, error)
	GetBasicInfoByUUID(ctx context.Context, uuid string) (*plugin_template_model.PluginTemplateBasicInfo, error)
	OnlineList(ctx context.Context, namespaceId int, uuid string) ([]*plugin_template_model.PluginTemplateOnlineItem, error)
	Online(ctx context.Context, namespaceId, operator int, uuid, clusterName string) (*frontend_model.Router, error)
	IsOnline(ctx context.Context, clusterId int, uuid string) (bool, error)
	Offline(ctx context.Context, namespaceId, operator int, uuid, clusterName string) error
	ResetOnline(ctx context.Context, namespaceId, clusterId int)
}
