package plugin

import (
	"context"
	plugin_model "github.com/eolinker/apinto-dashboard/modules/plugin/plugin-model"
)

type IClusterPluginService interface {
	GetList(ctx context.Context, namespaceID int, clusterName string) ([]*plugin_model.CluPluginListItem, error)
	GetPlugin(ctx context.Context, namespaceID int, clusterName string, pluginName string) (*plugin_model.ClusterPluginInfo, error)
	EditPlugin(ctx context.Context, namespaceID int, clusterName string, userID int, pluginName string, status int, config interface{}) error

	QueryHistory(ctx context.Context, namespaceId, pageNum, pageSize int, clusterName string) ([]*plugin_model.ClusterPluginHistory, int, error)
	PublishHistory(ctx context.Context, namespaceId, pageNum, pageSize int, clusterName string) ([]*plugin_model.ClusterPluginPublish, int, error)
	ToPublishes(ctx context.Context, namespaceId int, clusterName string) ([]*plugin_model.ClusterPluginToPublish, error)
	Publish(ctx context.Context, namespaceId, userId int, clusterName, versionName, desc, source string) error

	DeleteAll(ctx context.Context, clusterId int) error
	GetPublishVersion(ctx context.Context, clusterId int) (*plugin_model.ClusterPluginPublishVersion, error)
	ResetOnline(ctx context.Context, namespaceId, clusterId int)

	IsOnlineByName(ctx context.Context, namespaceID int, clusterName, pluginName string) (bool, error)
	IsDelete(ctx context.Context, namespaceID int, clusterName, pluginName string) (bool, error)
}
