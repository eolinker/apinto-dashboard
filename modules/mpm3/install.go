package mpm3

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/mpm3/model"
	"github.com/eolinker/apinto-dashboard/pm3"
)

type IInstallService interface {
	Install(ctx context.Context, cfg *pm3.PluginDefine, resource *model.PluginResources) error
	InstallInner(ctx context.Context, cfg *pm3.PluginDefine, resource *model.PluginResources, auto, canDisable bool) error
	Uninstall(ctx context.Context, pluginID string) error
	ClearInnerNotExits(ctx context.Context, ids []string) error

	CheckPluginInstalled(ctx context.Context, pluginID string) (bool, error)
}
