package mpm3

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/mpm3/model"
	"github.com/eolinker/apinto-dashboard/pm3"
)

type IAccessService interface {
	Save(ctx context.Context, plugin int, as []pm3.PAccess) error
	GetEnable(ctx context.Context) []*model.Access
	Clean(ctx context.Context)
}
