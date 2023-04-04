package middleware

import (
	"context"

	"github.com/eolinker/apinto-dashboard/modules/middleware/model"
)

type IMiddlewareService interface {
	CreateGroup(ctx context.Context, namespaceId int, operator int, uuid, prefix string, middlewares []string) error
	UpdateGroup(ctx context.Context, namespaceId int, operator int, uuid, prefix string, middlewares []string) error
	DeleteGroup(ctx context.Context, namespaceId int, operator int, uuid string) error
	GroupList(ctx context.Context, namespaceId int) ([]*model.MiddlewareGroup, error)
	GroupInfo(ctx context.Context, namespaceId int, uuid string) (*model.MiddlewareGroupInfo, error)
}
