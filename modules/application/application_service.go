package application

import (
	"context"
	driverInfo "github.com/eolinker/apinto-dashboard/driver-manager"
	"github.com/eolinker/apinto-dashboard/driver-manager/driver"
	application_dto2 "github.com/eolinker/apinto-dashboard/modules/application/application-dto"
	application_model2 "github.com/eolinker/apinto-dashboard/modules/application/application-model"
	"github.com/eolinker/apinto-dashboard/modules/online"
)

type IApplicationService interface {
	CreateApp(ctx context.Context, namespaceId, userId int, input *application_dto2.ApplicationInput) error
	UpdateApp(ctx context.Context, namespaceId, userId int, input *application_dto2.ApplicationInput) error
	DelApp(ctx context.Context, namespaceId, userId int, id string) error
	AppList(ctx context.Context, namespaceId, userId, pageNum, pageSize int, queryName string) ([]*application_model2.Application, int, error)
	AppListAll(ctx context.Context, namespaceId int) ([]*application_model2.Application, error)
	AppListFilter(ctx context.Context, namespaceId, pageNum, pageSize int, queryName string) ([]*application_model2.Application, int, error)
	AppListByUUIDS(ctx context.Context, namespaceId int, uuids []string) ([]*application_model2.Application, error)
	AppInfoDetails(ctx context.Context, namespaceId int, id string) (*application_model2.Application, error)
	AppInfo(ctx context.Context, namespaceId int, id string) (*application_model2.Application, error)
	Online(ctx context.Context, namespaceId, userId int, id, clusterName string) error
	Offline(ctx context.Context, namespaceId, userId int, id, clusterName string) error
	Disable(ctx context.Context, namespaceId, userId int, id, clusterName string, disable bool) error
	OnlineList(ctx context.Context, namespaceId int, id string) ([]*application_model2.ApplicationOnline, error)
	GetAppKeys(ctx context.Context, namespaceId int) ([]*application_model2.ApplicationKeys, error)
	GetAppVersion(ctx context.Context, appId int) (*application_model2.ApplicationVersion, error)
	online.IResetOnlineService
}

type IApplicationAuthService interface {
	GetList(ctx context.Context, namespaceId int, appId string) ([]*application_model2.ApplicationAuth, error)
	Create(ctx context.Context, namespaceId, userId int, appId string, input *application_dto2.ApplicationAuthInput) error
	Update(ctx context.Context, namespaceId, userId int, appId, uuid string, input *application_dto2.ApplicationAuthInput) error
	Delete(ctx context.Context, namespaceId, userId int, uuid string) error
	Info(ctx context.Context, namespaceId int, appId, uuid string) (*application_model2.ApplicationAuth, error)
	Online(ctx context.Context, namespaceId, userId, clusterId, applicationId int) ([]*application_model2.ApplicationAuth, error)
	Offline(ctx context.Context, clusterId, applicationId int) error
	IsUpdate(ctx context.Context, clusterId, applicationId int) (bool, error)
	GetListByApplicationId(ctx context.Context, applicationId int) ([]*application_model2.ApplicationAuth, error)
	GetDriversRender() []*driverInfo.DriverInfo
	GetDriver(driver string) driver.IAuthDriver
}
