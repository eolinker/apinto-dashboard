package middleware

import (
	"context"

	"github.com/eolinker/apinto-dashboard/modules/middleware/model"
)

type IMiddlewareService interface {
	Save(ctx context.Context, config []*model.Middleware) error
	Groups(ctx context.Context) ([]*model.Middleware, error)
}
