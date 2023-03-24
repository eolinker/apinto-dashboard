package plugin_model

import (
	"github.com/eolinker/apinto-dashboard/modules/plugin/plugin-entry"
	"time"
)

type CluPluginListItem struct {
	*plugin_entry.ClusterPlugin
	Operator     string
	Publish      int //1.未发布 2.已发布 3.缺失
	ChangeState  int //变更状态(增、改、删、无)
	ReleasedSort int //已发布的排序序号
	NowSort      int //现在的排序序号
	IsBuiltIn    bool
}

type ClusterPluginInfo struct {
	PluginName string
	Status     int
	Config     string
}

// ClusterPluginHistory 集群插件变更历史
type ClusterPluginHistory struct {
	*plugin_entry.ClusterPluginHistory
}

// ClusterPluginPublish 集群插件发布版本历史
type ClusterPluginPublish struct {
	Id         int
	Name       string
	OptType    int //1.发布 2.回滚
	Operator   string
	CreateTime time.Time
	Details    []*ClusterPluginPublishDetails
}

type ClusterPluginPublishDetails struct {
	Name         string
	OldValue     ClusterPluginHistoryConfig
	NewValue     ClusterPluginHistoryConfig
	ReleasedSort int
	NowSort      int
	OptType      int //1.新增 2.修改 3.删除
	CreateTime   time.Time
}

type ClusterPluginHistoryConfig struct {
	Status int
	Config string
}

// ClusterPluginToPublish 待发布的集群插件配置
type ClusterPluginToPublish struct {
	*plugin_entry.PluginToPublish
}

type ClusterPluginPublishVersion struct {
	*plugin_entry.ClusterPluginPublishVersion
}

type InnerPlugin struct {
	Id         string
	PluginName string
	Status     int
	Config     string
	Rely       string
}
