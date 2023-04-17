package module_plugin

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/dto"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/model"
)

type IModulePluginService interface {
	GetPlugins(ctx context.Context, groupID, searchName string) ([]*model.ModulePluginItem, error)
	GetPluginInfo(ctx context.Context, pluginUUID string) (*model.ModulePluginInfo, error)
	GetPluginGroups(ctx context.Context) ([]*model.PluginGroup, error)
	GetPluginEnableInfo(ctx context.Context, pluginUUID string) (*model.PluginEnableInfo, error)
	GetPluginEnableRender(ctx context.Context, pluginUUID string) (*model.PluginEnableRender, error)
	InstallPlugin(ctx context.Context, userID int, pluginYml *model.PluginYmlCfg, packageContent []byte) error
	UninstallPlugin(ctx context.Context, userID int, pluginID string) error
	EnablePlugin(ctx context.Context, userID int, pluginUUID string, enableInfo *dto.PluginEnableInfo) error
	DisablePlugin(ctx context.Context, userID int, pluginUUID string) error

	CheckPluginInstalled(ctx context.Context, pluginID string) (bool, error)
	CheckPluginISDeCompress(ctx context.Context, pluginDir string, pluginID string) error
	InstallInnerPlugin(ctx context.Context, pluginYml *model.InnerPluginYmlCfg) error
	UpdateInnerPlugin(ctx context.Context, pluginYml *model.InnerPluginYmlCfg) error
}

type IModulePlugin interface {
	//GetEnabledPlugins 获取已启用的插件信息列表
	GetEnabledPlugins(ctx context.Context) ([]*model.EnabledPlugin, error)
	//GetNavigationModules 获取导航接口所需要的模块列表
	GetNavigationModules(ctx context.Context) ([]*model.NavigationModuleInfo, error)
}
