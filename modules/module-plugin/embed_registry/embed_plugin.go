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
		pluginInfo, has := innerPluginsMap[p.PluginCfg.ID]
		if !has {
			// 插入安装记录
			err = service.InstallInnerPlugin(ctx, p.PluginCfg, p.Resources)
			if err != nil {
				return err
			}
			continue
		} else {
			//判断version有没改变，有则更新
			if pluginInfo.Version != p.PluginCfg.Version {
				err = service.UpdateInnerPlugin(ctx, pluginCfg)
				if err != nil {
					return err
				}
			}
			delete(innerPluginsMap, p.PluginCfg.ID)
		}
	}
	//TODO 遍历innerPluginsMap, 删除不存在的内置插件

	return nil
}
