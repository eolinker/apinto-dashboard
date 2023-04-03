package openapp

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/openapp/open-app-dto"
	"github.com/eolinker/apinto-dashboard/modules/openapp/open-app-model"
)

type IExternalApplicationService interface {
	AppList(ctx context.Context, namespaceId int) ([]*open_app_model.ExtAppListItem, error)
	AppInfo(ctx context.Context, namespaceId int, uuid string) (*open_app_model.ExternalAppInfo, error)
	CreateApp(ctx context.Context, namespaceId, userId int, input *open_app_dto.ExternalAppInfoInput) error
	UpdateApp(ctx context.Context, namespaceId, userId int, input *open_app_dto.ExternalAppInfoInput) error
	DelApp(ctx context.Context, namespaceId, userId int, uuid string) error
	Enable(ctx context.Context, namespaceId, userId int, uuid string) error
	Disable(ctx context.Context, namespaceId, userId int, uuid string) error
	FlushToken(ctx context.Context, namespaceId, userId int, uuid string) error

	CheckExtAPPToken(ctx context.Context, namespaceId int, token string) (int, error)

	UpdateExtAPPTags(ctx context.Context, namespaceId, appID int, label string) error
	GetExtAppName(ctx context.Context, id int) (string, error)
}
