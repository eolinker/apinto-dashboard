package service

import (
	context "context"
	"encoding/json"
	"errors"
	"fmt"
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
		enablePlugin := &model.EnabledPlugin{
			UUID:   p.UUID,
			Name:   p.Name,
			Driver: p.Driver,
			Config: enableCfg,
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
		//导航id不存在表示不需要在前端显示
		if IsInnerPlugin(module.Type) && module.Front == "" {
			continue
		}
		info := &model.NavigationModuleInfo{
			Name:       module.Name,
			Title:      module.Title,
			Type:       "built-in",
			Path:       module.Front,
			Navigation: module.Navigation,
		}
		//若模块为非内置模块
		if !IsInnerPlugin(module.Type) {
			info.Type = "outer"
			//TODO 临时处理
			if module.Front != "" {
				info.Path = module.Front
			} else {
				info.Path = fmt.Sprintf("/%s", module.Name)
			}
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

	info := &model.ModulePluginInfo{
		ModulePlugin: plugin,
		Enable:       true,
		CanDisable:   true,
		Uninstall:    false,
	}

	//根据类型判断是否能停用
	if IsPluginCanDisable(plugin.Type) {
		info.CanDisable = false
	}

	return info, nil

}
