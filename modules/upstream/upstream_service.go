package upstream

import (
	"context"

	frontend_model "github.com/eolinker/apinto-dashboard/modules/base/frontend-model"
	strategy_model "github.com/eolinker/apinto-dashboard/modules/strategy/strategy-model"
	upstream_model "github.com/eolinker/apinto-dashboard/modules/upstream/model"
	upstream_dto "github.com/eolinker/apinto-dashboard/modules/upstream/upstream-dto"
	upstream_entry2 "github.com/eolinker/apinto-dashboard/modules/upstream/upstream-entry"
)

type IService interface {
	GetServiceList(ctx context.Context, namespaceID int, searchName string, pageNum int, pageSize int) ([]*upstream_model.ServiceListItem, int, error)
	GetServiceListByNames(ctx context.Context, namespaceID int, names []string) ([]*upstream_model.ServiceListItem, error) //TODO service
	GetServiceListAll(ctx context.Context, namespaceID int, searchName string) ([]*upstream_model.ServiceListItem, error)  //TODO service
	GetServiceInfo(ctx context.Context, namespaceID int, serviceName string) (*upstream_model.ServiceInfo, error)
	CreateService(ctx context.Context, namespaceID, userId int, input *upstream_dto.ServiceInfo, variableList []string) (int, error) //TODO service
	UpdateService(ctx context.Context, namespaceID, userId int, input *upstream_dto.ServiceInfo, variableList []string) error
	DeleteService(ctx context.Context, namespaceID, userId int, serviceName string) error
	GetServiceEnum(ctx context.Context, namespaceID int, searchName string) ([]string, error)
	OnlineList(ctx context.Context, namespaceId int, serviceName string) ([]*upstream_model.ServiceOnline, error)
	OnlineService(ctx context.Context, namespaceId, operator int, serviceName, clusterName string) (*frontend_model.Router, error)
	OfflineService(ctx context.Context, namespaceId, operator int, serviceName, clusterName string) error
	GetServiceIDByName(ctx context.Context, namespaceId int, serviceName string) (int, error) //TODO service
	GetLatestServiceVersion(ctx context.Context, serviceID int) (*upstream_entry2.ServiceVersion, error)

	GetServiceSchemaInfo(ctx context.Context, serviceID int) (*upstream_entry2.Service, error)                                                      //TODO service
	IsOnline(ctx context.Context, clusterId, serviceId int) bool                                                                                    //TODO service
	GetServiceRemoteOptions(ctx context.Context, namespaceID, pageNum, pageSize int, keyword string) ([]*strategy_model.RemoteServices, int, error) //TODO service
	GetServiceRemoteByNames(ctx context.Context, namespaceID int, uuids []string) ([]*strategy_model.RemoteServices, error)                         //TODO service
	UpstreamCount(ctx context.Context, namespaceId int) (int64, error)                                                                              //TODO service

}
