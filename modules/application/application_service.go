package application

import (
	"context"
	"github.com/eolinker/apinto-dashboard/client/v1"
	driverInfo "github.com/eolinker/apinto-dashboard/driver"
	"github.com/eolinker/apinto-dashboard/modules/application/application-dto"
	"github.com/eolinker/apinto-dashboard/modules/application/application-model"
)

type IApplicationService interface {
	CreateApp(ctx context.Context, namespaceId, userId int, input *application_dto.ApplicationInput) error
	UpdateApp(ctx context.Context, namespaceId, userId int, input *application_dto.ApplicationInput) error
	DelApp(ctx context.Context, namespaceId, userId int, id string) error
	AppList(ctx context.Context, namespaceId, userId, pageNum, pageSize int, queryName string, clusters []string) ([]*application_model.ApplicationListItem, int, error)
	AppEnumList(ctx context.Context, namespaceId int) ([]*application_model.ApplicationBasicInfo, error)
	AllApp(ctx context.Context, namespaceId int) ([]*application_model.ApplicationBasicInfo, error)
	AppListByUUIDS(ctx context.Context, namespaceId int, uuids []string) ([]*application_model.ApplicationBasicInfo, error)
	AppBasicInfo(ctx context.Context, namespaceId int, uuid string) (*application_model.ApplicationBasicInfo, error)
	AppInfoDetails(ctx context.Context, namespaceId int, id string) (*application_model.ApplicationInfo, error)
	AppInfo(ctx context.Context, namespaceId int, id string) (*application_model.ApplicationEntire, error)
	Online(ctx context.Context, namespaceId, userId int, id string, clusterNames []string) error
	Offline(ctx context.Context, namespaceId, userId int, id string, clusterNames []string) error
	OnlineInfo(ctx context.Context, namespaceId int, uuid string) (*application_model.ApplicationBasicInfo, []*application_model.AppCluster, error)
	GetAppKeys(ctx context.Context, namespaceId int) ([]*application_model.ApplicationKeys, error)
	GetAppRemoteOptions(ctx context.Context, namespaceId, pageNum, pageSize int, keyword string) ([]any, error)
	//online.IResetOnlineService

	GetAuthList(ctx context.Context, namespaceId int, appId string) ([]*application_model.AppAuthItem, error)
	CreateAuth(ctx context.Context, namespaceId, userId int, appId string, input *application_dto.ApplicationAuthInput) error
	UpdateAuth(ctx context.Context, namespaceId, userId int, appId, uuid string, input *application_dto.ApplicationAuthInput) error
	DeleteAuth(ctx context.Context, namespaceId, userId int, uuid string) error
	AuthInfo(ctx context.Context, namespaceId int, appId, uuid string) (*application_model.ApplicationAuth, error)
	AuthDetails(ctx context.Context, namespaceId int, appId, uuid string) ([]application_model.AuthDetailItem, error)
	OnlineAuth(ctx context.Context, applicationId int) ([]*application_model.ApplicationAuth, error)
	GetDriversRender() []*driverInfo.DriverInfo
	GetDriver(driver string) IAuthDriver
}

type IAuthDriverManager interface {
	driverInfo.IDriverManager[IAuthDriver]
	List() []*driverInfo.DriverInfo
}

type IAuthDriver interface {
	Render() string
	CheckInput(config []byte) error
	//GetAuthListInfo 获取健全列表展示需要用的参数信息
	GetAuthListInfo(config []byte) string
	GetCfgDetails(config []byte) []application_model.AuthDetailItem
	ToApinto(expire int64, position string, tokenName string, config []byte, hideCredential bool) v1.ApplicationAuth
}
