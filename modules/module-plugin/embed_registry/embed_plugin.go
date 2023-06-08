package embed_registry

import (
	"context"
	"github.com/eolinker/apinto-dashboard/common"
	module_plugin "github.com/eolinker/apinto-dashboard/modules/module-plugin"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/model"
	"github.com/eolinker/eosc/common/bean"
)

var (
	embedPlugins []*model.EmbedPluginCfg
)

func RegisterEmbedPlugin(ps ...*model.EmbedPluginCfg) {
	embedPlugins = append(embedPlugins, ps...)
}

func InitEmbedPlugins() error {
	var service module_plugin.IModulePluginService
	bean.Autowired(&service)
	ctx := context.Background()

	innerPlugins, err := service.GetInnerPluginList(ctx)
	if err != nil {
		return err
	}
	innerPluginsMap := common.SliceToMap(innerPlugins, func(t *model.ModulePluginInfo) string {
		return t.UUID
	})
	for _, p := range embedPlugins {
		//TODO 校验内置插件
		epCfg := p.PluginCfg
		pluginInfo, has := innerPluginsMap[epCfg.ID]
		cfg := &model.PluginCfg{
			Version:    epCfg.Version,
			Navigation: epCfg.Navigation,
			GroupID:    epCfg.GroupID,
			Resume:     epCfg.Resume,
			Type:       epCfg.Type,
			Define:     epCfg.Define,
		}
		if !has {
			// 插入安装记录
			err = service.InstallInnerPlugin(ctx, epCfg.ID, epCfg.Name, epCfg.CName, epCfg.Driver, epCfg.Icon, epCfg.Auto, epCfg.IsCanDisable,
				epCfg.IsCanUninstall, epCfg.VisibleInNavigation, epCfg.VisibleInMarket, cfg, p.Resources)
			if err != nil {
				return err
			}
			continue
		} else {
			//判断version有没改变，有则更新
			if pluginInfo.Version != p.PluginCfg.Version {
				err = service.UpdateInnerPlugin(ctx, epCfg.ID, epCfg.Name, epCfg.CName, epCfg.Driver, epCfg.Icon, epCfg.IsCanDisable,
					epCfg.IsCanUninstall, epCfg.VisibleInNavigation, epCfg.VisibleInMarket, cfg, p.Resources)
				if err != nil {
					return err
				}
			}
			delete(innerPluginsMap, p.PluginCfg.ID)
		}
	}

	//遍历innerPluginsMap, 删除不存在的内置插件
	deleteIds := make([]int, 0, len(innerPluginsMap))
	for _, v := range innerPluginsMap {
		deleteIds = append(deleteIds, v.Id)
	}

	return service.DeleteInnerByIds(ctx, deleteIds...)
}
