package service

import (
	context "context"
	"encoding/json"
	"errors"
	locker_service "github.com/eolinker/apinto-dashboard/modules/base/locker-service"
	"github.com/eolinker/apinto-dashboard/modules/group"
	module_plugin "github.com/eolinker/apinto-dashboard/modules/module-plugin"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/model"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/store"
	"github.com/eolinker/eosc/common/bean"
	"gorm.io/gorm"
)

type modulePlugin struct {
	pluginStore        store.IModulePluginStore
	pluginEnableStore  store.IModulePluginEnableStore
	pluginPackageStore store.IModulePluginPackageStore

	commonGroup group.ICommonGroupService
	lockService locker_service.IAsynLockService
}

func newModulePlugin() module_plugin.IModulePlugin {

	s := &modulePlugin{}
	bean.Autowired(&s.pluginStore)
	bean.Autowired(&s.pluginEnableStore)
	bean.Autowired(&s.pluginPackageStore)

	bean.Autowired(&s.commonGroup)
	bean.Autowired(&s.lockService)
	return s
}

func (m *modulePlugin) GetEnabledPlugins(ctx context.Context) ([]*model.EnabledPlugin, error) {
	plugins, err := m.pluginStore.GetEnabledPlugins(ctx)
	if err != nil {
		return nil, err
	}
	enablePlugins := make([]*model.EnabledPlugin, 0, len(plugins))
	for _, p := range plugins {
		enableCfg := new(model.PluginEnableCfg)
		_ = json.Unmarshal(p.Config, enableCfg)
		enableCfgMap := enabledCfgListToMap(enableCfg)
		enablePlugin := &model.EnabledPlugin{
			UUID:   p.UUID,
			Name:   p.Name,
			Driver: p.Driver,
			Config: enableCfgMap,
			Define: nil,
		}
		define := new(model.PluginDefine)
		_ = json.Unmarshal(p.Details, define)
		enablePlugin.Define = define

		enablePlugins = append(enablePlugins, enablePlugin)
	}
	return enablePlugins, nil
}

func (m *modulePlugin) GetNavigationModules(ctx context.Context) ([]*model.NavigationModuleInfo, error) {
	moduleInfos, err := m.pluginStore.GetNavigationModules(ctx)
	if err != nil {
		return nil, err
	}

	list := make([]*model.NavigationModuleInfo, 0, len(moduleInfos))
	for _, module := range moduleInfos {
		//导航不存在表示不需要在前端显示
		if !module.IsPluginVisible {
			continue
		}
		info := &model.NavigationModuleInfo{
			Name:       module.Name,
			Title:      module.Title,
			Path:       module.Frontend,
			Navigation: module.Navigation,
		}

		list = append(list, info)
	}
	return list, nil
}

func (m *modulePlugin) GetEnabledPluginByModuleName(ctx context.Context, moduleName string) (*model.ModulePluginInfo, error) {
	enableInfo, err := m.pluginEnableStore.GetEnabledPluginByName(ctx, moduleName)
	if err != nil && err != gorm.ErrRecordNotFound {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("插件未启用")
		}
		return nil, err
	}

	plugin, err := m.pluginStore.Get(ctx, enableInfo.Id)
	if err != nil {
		return nil, err
	}

	return &model.ModulePluginInfo{
		ModulePlugin: plugin,
		Enable:       true,
		CanDisable:   enableInfo.IsCanDisable,
		Uninstall:    enableInfo.IsCanUninstall,
	}, nil

}
