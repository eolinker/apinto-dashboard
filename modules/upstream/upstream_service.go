package upstream

import (
	"context"
	"github.com/eolinker/apinto-dashboard/dto/service-dto"
	"github.com/eolinker/apinto-dashboard/entry/upstream-entry"
	"github.com/eolinker/apinto-dashboard/model"
	upstream_model "github.com/eolinker/apinto-dashboard/modules/upstream/model"
)

type IService interface {
	GetServiceList(ctx context.Context, namespaceID int, searchName string, pageNum int, pageSize int) ([]*upstream_model.ServiceListItem, int, error)
	GetServiceListByNames(ctx context.Context, namespaceID int, names []string) ([]*upstream_model.ServiceListItem, error)
	GetServiceListAll(ctx context.Context, namespaceID int) ([]*upstream_model.ServiceListItem, error)
	GetServiceInfo(ctx context.Context, namespaceID int, serviceName string) (*upstream_model.ServiceInfo, error)
	CreateService(ctx context.Context, namespaceID, userId int, input *service_dto.ServiceInfo, variableList []string) (int, error)
	UpdateService(ctx context.Context, namespaceID, userId int, input *service_dto.ServiceInfo, variableList []string) error
	DeleteService(ctx context.Context, namespaceID, userId int, serviceName string) error
	GetServiceEnum(ctx context.Context, namespaceID int, searchName string) ([]string, error)
	OnlineList(ctx context.Context, namespaceId int, serviceName string) ([]*upstream_model.ServiceOnline, error)
	OnlineService(ctx context.Context, namespaceId, operator int, serviceName, clusterName string) (*model.Router, error)
	OfflineService(ctx context.Context, namespaceId, operator int, serviceName, clusterName string) error
	GetServiceIDByName(ctx context.Context, namespaceId int, serviceName string) (int, error)
	GetLatestServiceVersion(ctx context.Context, serviceID int) (*upstream_entry.ServiceVersion, error)
	GetServiceSchemaInfo(ctx context.Context, serviceID int) (*upstream_entry.Service, error)
	IsOnline(ctx context.Context, clusterId, serviceId int) bool
	GetServiceRemoteOptions(ctx context.Context, namespaceID, pageNum, pageSize int, keyword string) ([]*model.RemoteServices, int, error)
	GetServiceRemoteByNames(ctx context.Context, namespaceID int, uuids []string) ([]*model.RemoteServices, error)
	ResetOnline(ctx context.Context, namespaceId, clusterId int)
}
