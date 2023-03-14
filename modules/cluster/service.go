package cluster

import (
	"context"
	v1 "github.com/eolinker/apinto-dashboard/client/v1"
	"github.com/eolinker/apinto-dashboard/modules/cluster/cluster-dto"
	"github.com/eolinker/apinto-dashboard/modules/cluster/cluster-entry"
	cluster_model2 "github.com/eolinker/apinto-dashboard/modules/cluster/cluster-model"
	"github.com/eolinker/apinto-dashboard/modules/online"
)

type IApintoClient interface {
	GetClient(ctx context.Context, clusterId int) (v1.IClient, error)
	SetClient(namespace, clusterId int)
	InitClustersGlobalPlugin(ctx context.Context) error
	InitGlobalPlugin(clusterAddr string, nodesAdminAddr []string) error
}
type IClusterCertificateService interface {
	Insert(ctx context.Context, operator, namespaceId int, clusterName, key, pem string) error
	Update(ctx context.Context, operator, namespaceId, certificateId int, clusterName, key, pem string) error
	QueryList(ctx context.Context, namespaceId int, clusterName string) ([]*cluster_model2.ClusterCertificate, error)
	DeleteById(ctx context.Context, namespaceId int, clusterName string, id int) error
}

type IClusterService interface {
	GetAllCluster(ctx context.Context) ([]*cluster_model2.Cluster, error)
	CheckByNamespaceByName(ctx context.Context, namespaceId int, name string) (int, error)
	GetByClusterId(ctx context.Context, clusterId int) (*cluster_model2.Cluster, error)
	GetByNamespaceByName(ctx context.Context, namespaceId int, name string) (*cluster_model2.Cluster, error)
	GetByNamespaceId(ctx context.Context, namespaceId int) ([]*cluster_model2.Cluster, error)
	GetByNames(ctx context.Context, namespaceId int, names []string) ([]*cluster_model2.Cluster, error)
	Insert(ctx context.Context, namespaceId, userId int, clusterInput *cluster_dto.ClusterInput) error
	QueryByNamespaceId(ctx context.Context, namespaceId int, clusterName string) (*cluster_model2.Cluster, error)
	QueryListByNamespaceId(ctx context.Context, namespaceId int) ([]*cluster_model2.Cluster, error)
	DeleteByNamespaceIdByName(ctx context.Context, namespaceId, userId int, name string) error
	UpdateDesc(ctx context.Context, namespaceId, userId int, name, desc string) error
	UpdateAddr(ctx context.Context, userId, clusterId int, addr, uuid string) error
}
type IClusterConfigService interface {
	Get(ctx context.Context, namespaceId int, clusterName, configType string) (interface{}, error)
	Edit(ctx context.Context, namespaceId, operator int, clusterName, configType string, config []byte) error
	Enable(ctx context.Context, namespaceId, operator int, clusterName, configType string) error
	Disable(ctx context.Context, namespaceId, operator int, clusterName, configType string) error

	IsConfigTypeExist(configType string) bool
	CheckInput(configType string, config []byte) error
	FormatOutput(configType string, operator string, config *cluster_entry.ClusterConfig) interface{}
	ToApinto(client v1.IClient, name, configType string, config []byte) error
	OfflineApinto(client v1.IClient, name, configType string) error
	online.IResetOnlineService
}

type IClusterNodeService interface {
	QueryList(ctx context.Context, namespaceId int, clusterName string) ([]*cluster_model2.ClusterNode, bool, error)
	QueryByClusterIds(ctx context.Context, clusterIds ...int) ([]*cluster_model2.ClusterNode, error)
	Reset(ctx context.Context, namespaceId, userId int, clusterName, clusterAddr, source string) error
	Update(ctx context.Context, namespaceId int, clusterName string) error
	NodeRepeatContrast(ctx context.Context, namespaceId, clusterId int, newList []*cluster_model2.ClusterNode) error
	Insert(ctx context.Context, nodes []*cluster_model2.ClusterNode) error
	GetNodesByUrl(addr string) ([]*cluster_model2.ClusterNode, error)
	GetClusterInfo(addr string) (*v1.ClusterInfo, error)
}
