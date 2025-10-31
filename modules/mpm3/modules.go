package mpm3

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/mpm3/model"
	"github.com/eolinker/apinto-dashboard/pm3"
)

type IModuleService interface {
	Save(ctx context.Context, plugin int, module []pm3.PModule) error
	GetEnable(ctx context.Context) []*model.Module
	Clean(ctx context.Context)
}
