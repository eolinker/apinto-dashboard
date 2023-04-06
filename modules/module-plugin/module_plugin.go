package module_plugin

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/dto"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/model"
)

type IModulePluginService interface {
	GetPlugins(ctx context.Context, groupUUID, searchName string) ([]*model.ModulePluginItem, []*model.PluginGroup, error)
	GetPluginInfo(ctx context.Context, pluginUUID string) (*model.ModulePluginInfo, error)
	GetPluginGroupsEnum(ctx context.Context) ([]*model.PluginGroup, error)
	GetPluginEnableInfo(ctx context.Context, pluginUUID string) (*model.PluginEnableInfo, *model.PluginEnableRender, error)
	InstallPlugin(ctx context.Context, groupName string, packageContent []byte) error
	EnablePlugin(ctx context.Context, pluginUUID string, enableInfo *dto.PluginEnableInfo) error
	DisablePlugin(ctx context.Context, pluginUUID string) error
}
