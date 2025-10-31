package mpm3

import (
	"context"
	"github.com/eolinker/apinto-dashboard/pm3"
)

type IFrontendService interface {
	Save(ctx context.Context, plugin int, content []pm3.PFrontend) error
	GetEnable(ctx context.Context) []pm3.PFrontend
	Clean(ctx context.Context)
}
