package module_plugin

import (
	"context"
	"github.com/eolinker/apinto-dashboard/cache"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/dto"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/entry"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/model"
)

type IModulePluginService interface {
	GetPlugins(ctx context.Context, groupID, searchName string) ([]*model.ModulePluginItem, error)
	GetPluginInfo(ctx context.Context, pluginUUID string) (*model.ModulePluginInfo, error)
	GetPluginGroups(ctx context.Context) ([]*model.PluginGroup, error)
	GetInnerPluginList(ctx context.Context) ([]*model.ModulePluginInfo, error)
	GetPluginEnableInfo(ctx context.Context, pluginUUID string) (*model.PluginEnableInfo, error)
	GetPluginEnableRender(ctx context.Context, pluginUUID string) (*model.PluginEnableRender, error)

	InstallPlugin(ctx context.Context, userID int, id, name, cname, driver, icon string, cfg *model.PluginCfg, resources *model.PluginResources) error
	InstallInnerPlugin(ctx context.Context, id, name, cname, driver, icon string, isEnable, isCanDisable, isCanUninstall, visibleInNavigation, visibleInMarket bool, cfg *model.PluginCfg, resources *model.EmbedPluginResources) error
	Install(ctx context.Context, userID int, id, name, cname, driver, icon string, isCanDisable, isCanUninstall, isInner, visibleInNavigation, visibleInMarket bool, cfg *model.PluginCfg, resources *model.PluginResources) error
	UpdateInnerPlugin(ctx context.Context, id, name, cname, driver, icon string, isCanDisable, isCanUninstall, visibleInNavigation, visibleInMarket bool, cfg *model.PluginCfg, resources *model.EmbedPluginResources) error
	UninstallPlugin(ctx context.Context, pluginID string) error
	DeleteInnerByIds(ctx context.Context, ids ...int) error
	EnablePlugin(ctx context.Context, userID int, pluginUUID string, enableInfo *dto.PluginEnableInfo) error
	DisablePlugin(ctx context.Context, userID int, pluginUUID string) error

	CheckPluginInstalled(ctx context.Context, pluginID string) (bool, error)
	CheckExternPluginInCache(ctx context.Context, pluginID string) error
}

type IModulePlugin interface {
	//GetEnabledPlugins 获取已启用的插件信息列表
	GetEnabledPlugins(ctx context.Context) ([]*model.EnabledPlugin, error)
	//GetNavigationModules 获取导航接口所需要的模块列表
	GetNavigationModules(ctx context.Context) ([]*model.NavigationModuleInfo, error)
	//GetEnabledPluginByModuleName 根据已启用插件的模块名获取插件信息
	GetEnabledPluginByModuleName(ctx context.Context, moduleName string) (*model.ModulePluginInfo, error)
}

type INavigationModulesCache interface {
	cache.IRedisCacheNoKey[entry.EnabledModule]
}
