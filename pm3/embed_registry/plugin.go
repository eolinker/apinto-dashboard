package embed_registry

import (
	"context"
	"github.com/eolinker/apinto-dashboard/custom"
	"github.com/eolinker/apinto-dashboard/modules/mpm3"
	"github.com/eolinker/eosc/common/bean"
)

var (
	embedPlugins []*PluginCfg
	service      mpm3.IInstallService
)

func registerEmbedPlugin(ps ...*PluginCfg) {
	embedPlugins = append(embedPlugins, ps...)
}

func init() {
	bean.Autowired(&service)
}
func InitEmbedPlugins() error {
	bean.AddInitializingBeanFunc(func() {
		ctx := context.Background()

		defer func() {
			embedPlugins = nil
		}()
		ids := make([]string, 0, len(embedPlugins))
		for _, p := range embedPlugins {
			if custom.IgnorePlugin(p.define.Id) {
				continue
			}
			err := service.InstallInner(ctx, p.define, p.resources, p.auto, p.isDisable)
			if err != nil {
				panic(err)
				return
			}
			ids = append(ids, p.define.Id)

		}

		service.ClearInnerNotExits(ctx, ids)
	})
	bean.Autowired(&service)

	return nil
}
