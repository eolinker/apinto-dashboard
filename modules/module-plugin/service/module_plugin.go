package service

import (
	context "context"
	"encoding/json"
	locker_service "github.com/eolinker/apinto-dashboard/modules/base/locker-service"
	"github.com/eolinker/apinto-dashboard/modules/group"
	module_plugin "github.com/eolinker/apinto-dashboard/modules/module-plugin"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/model"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/store"
	"github.com/eolinker/eosc/common/bean"
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
		switch p.Driver {
		case "remote":
			remote := new(model.RemoteDefine)
			_ = json.Unmarshal(p.Define, remote)
			enablePlugin.Define = remote
		case "local":
			local := new(model.LocalDefine)
			_ = json.Unmarshal(p.Define, local)
			enablePlugin.Define = local
		case "profession":
			profession := new(model.ProfessionDefine)
			_ = json.Unmarshal(p.Define, profession)
			enablePlugin.Define = profession
		default:

		}

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
		//路由为空的不需要返回
		if module.Front == "" {
			continue
		}
		info := &model.NavigationModuleInfo{
			Name:       module.Name,
			Title:      module.Title,
			Type:       "outer",
			Path:       "",
			Navigation: module.Navigation,
		}
		//若模块为内置模块
		if module.Type == 0 || module.Type == 1 {
			info.Type = "built-in"
			info.Path = module.Front
		}
		list = append(list, info)
	}
	return list, nil
}
