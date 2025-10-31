package mpm3

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/mpm3/model"
)

type IPluginService interface {
	GetPlugin(ctx context.Context, uuid string) (*model.PluginInfo, error)
	GetEnabled(ctx context.Context) ([]*model.PluginConfig, error)
	GetEnableRender(ctx context.Context, pluginUUID string) (*model.PluginEnableRender, error)

	GetPluginEnableInfo(ctx context.Context, pluginUUID string) (*model.PluginEnableCfg, error)

	Search(ctx context.Context, groupUUID, searchName string) ([]*model.Plugin, error)
	GetGroups(ctx context.Context) ([]*model.PluginGroup, error)

	EnablePlugin(ctx context.Context, userID int, pluginUUID string, enableInfo *model.PluginEnableCfg) error
	DisablePlugin(ctx context.Context, userID int, pluginUUID string) error
}
