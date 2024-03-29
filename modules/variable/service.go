package variable

import (
	"context"

	quote_entry "github.com/eolinker/apinto-dashboard/modules/base/quote-entry"
	"github.com/eolinker/apinto-dashboard/modules/cluster/cluster-dto"
	"github.com/eolinker/apinto-dashboard/modules/variable/variable-model"
)

type IGlobalVariableService interface {
	List(ctx context.Context, pageNum, pageSize, namespace int, key string, status int) ([]*variable_model.GlobalVariableListItem, int, error)
	GetInfo(ctx context.Context, namespaceID int, key string) ([]*variable_model.GlobalVariableDetails, error)
	Create(ctx context.Context, namespaceID, userID int, key, desc string) (int, error)
	Delete(ctx context.Context, namespaceID, userID int, key string) error
	GetByKeys(ctx context.Context, namespaceId int, keys []string) ([]*variable_model.GlobalVariable, error)
	GetById(ctx context.Context, id int) (*variable_model.GlobalVariable, error)
	//QuoteVariables 引用环境变量
	QuoteVariables(ctx context.Context, namespaceID, sourceID int, quoteType quote_entry.QuoteKindType, variableKeys []string) error
	//CheckQuotedVariablesOnline 检查Source引用的环境变量是否在目标集群已上线
	CheckQuotedVariablesOnline(ctx context.Context, clusterID int, clusterName string, sourceID int, quoteType quote_entry.QuoteKindType) error
	//DeleteVariableQuote 删除source的环境变量引用关系
	DeleteVariableQuote(ctx context.Context, sourceID int, quoteType quote_entry.QuoteKindType) error
}

type IClusterVariableService interface {
	GetList(ctx context.Context, namespaceID int, clusterName string) ([]*variable_model.ClusterVariableListItem, error)
	Create(ctx context.Context, namespaceID int, clusterName string, userID int, key, value, desc string) error
	Update(ctx context.Context, namespaceID int, clusterName string, userID int, key, value string) error
	Delete(ctx context.Context, namespaceID int, clusterName string, userID int, key string) error
	DeleteAll(ctx context.Context, namespaceID int, clusterId, userID int) error
	SyncConf(ctx context.Context, namespaceId, userId int, clusterName string, conf *cluster_dto.SyncConf) error
	QueryHistory(ctx context.Context, namespaceId, pageNum, pageSize int, clusterName string) ([]*variable_model.ClusterVariableHistory, int, error)
	ToPublishs(ctx context.Context, namespaceId int, clusterName string) ([]*variable_model.VariableToPublish, error)
	Publish(ctx context.Context, namespaceId, userId int, clusterName, versionName, desc, source string) error
	GetSyncConf(ctx context.Context, namespaceId int, clusterName string) (*variable_model.ClustersVariables, error)
	PublishHistory(ctx context.Context, namespaceId, pageNum, pageSize int, clusterName string) ([]*variable_model.VariablePublish, int, error)
	GetPublishVersion(ctx context.Context, clusterId int) (*variable_model.VariablePublishVersion, error)
	//online.IResetOnlineService
}
