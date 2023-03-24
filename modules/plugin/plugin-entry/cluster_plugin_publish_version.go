package plugin_entry

import (
	"fmt"
	"time"
)

// ClusterPluginPublishVersion 集群绑定的插件版本信息表
type ClusterPluginPublishVersion struct {
	Id                   int                       `json:"id"`
	ClusterId            int                       `json:"cluster_id"`
	NamespaceId          int                       `json:"namespace_id"`
	PublishedPluginsList []*CluPluginPublishConfig `json:"published_plugins_list"`
	Desc                 string                    `json:"desc"`
	Operator             int                       `json:"operator"`
	CreateTime           time.Time                 `json:"create_time"`
}

func (c *ClusterPluginPublishVersion) SetVersionId(id int) {
	c.Id = id
}

type CluPluginPublishConfig struct {
	ClusterPlugin *ClusterPlugin `json:"cluster_plugin"`
	Sort          int            `json:"sort"`
	Extended      string         `json:"extended"`
	Rely          string         `json:"rely"`
}

func (c *CluPluginPublishConfig) GetKey() string {
	return c.ClusterPlugin.PluginName
}

func (c *CluPluginPublishConfig) Values() []string {
	return []string{fmt.Sprintf("%d-%d-%s", c.Sort, c.ClusterPlugin.Status, c.ClusterPlugin.Config)}
}
