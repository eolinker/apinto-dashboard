package service

import (
	context "context"
	"encoding/json"
	"fmt"
	locker_service "github.com/eolinker/apinto-dashboard/modules/base/locker-service"
	"github.com/eolinker/apinto-dashboard/modules/group"
	module_plugin "github.com/eolinker/apinto-dashboard/modules/module-plugin"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/model"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/store"
	"github.com/eolinker/apinto-dashboard/modules/navigation"
	"github.com/eolinker/eosc/common/bean"
)

type modulePlugin struct {
	pluginStore        store.IModulePluginStore
	pluginEnableStore  store.IModulePluginEnableStore
	pluginPackageStore store.IModulePluginPackageStore

	commonGroup       group.ICommonGroupService
	navigationService navigation.INavigationService
	lockService       locker_service.IAsynLockService
}

func newModulePlugin() module_plugin.IModulePlugin {

	s := &modulePlugin{}
	bean.Autowired(&s.pluginStore)
	bean.Autowired(&s.pluginEnableStore)
	bean.Autowired(&s.pluginPackageStore)

	bean.Autowired(&s.commonGroup)
	bean.Autowired(&s.navigationService)
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

func (m *modulePlugin) GetMiddlewareList(ctx context.Context) ([]*model.MiddlewareItem, error) {
	plugins, err := m.pluginStore.GetEnabledPlugins(ctx)
	if err != nil {
		return nil, err
	}
	middlewares := make([]*model.MiddlewareItem, 0, len(plugins))
	for _, p := range plugins {
		switch p.Driver {
		case "remote", "profession":
		case "local":
			local := new(model.LocalDefine)
			_ = json.Unmarshal(p.Define, local)
			for _, l := range local.Middleware {
				middlewares = append(middlewares, &model.MiddlewareItem{
					Name: fmt.Sprintf("%s.%s", p.Name, l.Name),
				})
			}
			//内置插件
		default:
			inner := new(model.InnerDefine)
			_ = json.Unmarshal(p.Define, inner)
			for _, i := range inner.Main.Middleware {
				middlewares = append(middlewares, &model.MiddlewareItem{
					Name: fmt.Sprintf("%s.%s", p.Name, i),
					Desc: "",
				})
			}
		}

	}
	return middlewares, nil
}
