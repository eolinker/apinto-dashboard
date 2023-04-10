package module_plugin

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/dto"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/model"
)

type IModulePluginService interface {
	GetPlugins(ctx context.Context, groupUUID, searchName string) ([]*model.ModulePluginItem, error)
	GetPluginInfo(ctx context.Context, pluginUUID string) (*model.ModulePluginInfo, error)
	GetPluginGroups(ctx context.Context) ([]*model.PluginGroup, error)
	GetPluginEnableInfo(ctx context.Context, pluginUUID string) (*model.PluginEnableInfo, error)
	GetPluginEnableRender(ctx context.Context, pluginUUID string) (*model.PluginEnableRender, error)
	InstallPlugin(ctx context.Context, userID int, groupName string, pluginYml *model.PluginYmlCfg, packageContent []byte) error
	EnablePlugin(ctx context.Context, userID int, pluginUUID string, enableInfo *dto.PluginEnableInfo) error
	DisablePlugin(ctx context.Context, userID int, pluginUUID string) error

	GetEnablePluginsByNavigation(ctx context.Context, navigationID int) ([]*model.NavigationEnabledPlugin, error)
	CheckPluginInstalled(ctx context.Context, pluginID string) (bool, error)
}

type IModulePlugin interface {
	InstallInnerPlugin(ctx context.Context, pluginYml *model.InnerPluginYmlCfg) error
	//GetEnabledPlugins 获取已启用的插件信息列表
	GetEnabledPlugins(ctx context.Context) ([]*model.EnabledPlugin, error)
	//GetMiddlewareList 获取拦截器列表
	GetMiddlewareList(ctx context.Context) ([]*model.MiddlewareItem, error)
}
