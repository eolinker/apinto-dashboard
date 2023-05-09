package discovery

import (
	"context"
	driver_manager "github.com/eolinker/apinto-dashboard/driver"
	"github.com/eolinker/apinto-dashboard/modules/base/frontend-model"
	"github.com/eolinker/apinto-dashboard/modules/discovery/discover-dto"
	"github.com/eolinker/apinto-dashboard/modules/discovery/discovery-model"
	"github.com/eolinker/apinto-dashboard/modules/upstream"
)

type IDiscoveryService interface {
	GetDiscoveryList(ctx context.Context, namespaceID int, searchName string) ([]*discovery_model.DiscoveryListItem, error)
	GetDiscoveryVersionInfo(ctx context.Context, namespaceID int, discoveryName string) (*discovery_model.DiscoveryInfo, error)
	CreateDiscovery(ctx context.Context, namespaceID int, userID int, input *discover_dto.DiscoveryInfoProxy) error
	UpdateDiscovery(ctx context.Context, namespaceID int, userID int, input *discover_dto.DiscoveryInfoProxy) error
	DeleteDiscovery(ctx context.Context, namespaceId, userId int, discoveryName string) error

	OnlineList(ctx context.Context, namespaceId int, discoveryName string) ([]*discovery_model.DiscoveryOnline, error)
	OnlineDiscovery(ctx context.Context, namespaceId, operator int, discoveryName, clusterName string) (*frontend_model.Router, error)
	OfflineDiscovery(ctx context.Context, namespaceId, operator int, discoveryName, clusterName string) error

	GetDiscoveryName(ctx context.Context, discoveryID int) (string, error)
	GetDiscoveryInfoByID(ctx context.Context, discoveryID int) (*discovery_model.DiscoveryListItem, error)
	GetDiscoveryID(ctx context.Context, namespaceID int, discoveryName string) (int, error)
	GetDiscoveryEnum(ctx context.Context, namespaceID int) ([]*discovery_model.DiscoveryEnum, error)
	GetDriversRender() []*driver_manager.DriverInfo
	GetLatestDiscoveryVersion(ctx context.Context, discoveryID int) (*discovery_model.DiscoveryVersion, error)
	IsOnline(ctx context.Context, clusterId, discoveryId int) bool

	//通过服务名获取配置上游服务时所需要的discoveryDriver
	GetServiceDiscoveryDriver(ctx context.Context, namespaceID int, discoveryName string) (int, string, upstream.IServiceDriver, error)
	GetServiceDiscoveryDriverByID(ctx context.Context, discoveryID int) (string, string, upstream.IServiceDriver, error)
	//online.IResetOnlineService
}
