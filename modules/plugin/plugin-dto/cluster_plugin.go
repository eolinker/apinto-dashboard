package plugin_dto

import "github.com/eolinker/apinto-dashboard/enum"

type CluPluginListItem struct {
	Name         string               `json:"name"`
	Publish      enum.PublishType     `json:"publish"`
	Config       string               `json:"config"`
	Status       enum.PluginStateType `json:"status"`
	ChangeStatus enum.ChangeOptType   `json:"change_status"`
	ReleasedSort int                  `json:"released_sort"`
	NowSort      int                  `json:"now_sort"`
	IsBuiltIn    bool                 `json:"is_builtin"`
	Operator     string               `json:"operator"`
	CreateTime   string               `json:"create_time"`
	UpdateTime   string               `json:"update_time"`
}

type ClusterPluginInfoInput struct {
	PluginName string      `json:"name"`
	Status     string      `json:"status"`
	Config     interface{} `json:"config"`
}

type ClusterPluginInfo struct {
	PluginName string               `json:"name"`
	Status     enum.PluginStateType `json:"status"`
	Config     interface{}          `json:"config"`
}

// ClusterPluginHistoryItem 集群插件变更历史项
type ClusterPluginHistoryItem struct {
	Name       string              `json:"name"`
	OldConfig  ClusterPluginConfig `json:"old_config"`
	NewConfig  ClusterPluginConfig `json:"new_config"`
	CreateTime string              `json:"create_time"`
	OptType    enum.ChangeOptType  `json:"opt_type"`
}

type ClusterPluginConfig struct {
	Status enum.PluginStateType `json:"status"`
	Config string               `json:"config"`
}

// ClusterPluginPublishItem 集群插件发布历史项
type ClusterPluginPublishItem struct {
	Id         int                            `json:"id"`
	Name       string                         `json:"name"`
	OptType    enum.PublishOptType            `json:"opt_type"`
	Operator   string                         `json:"operator"`
	CreateTime string                         `json:"create_time"`
	Details    []*ClusterPluginPublishDetails `json:"details"`
}

type ClusterPluginPublishDetails struct {
	Name         string              `json:"name"`
	OldConfig    ClusterPluginConfig `json:"old_config"`
	NewConfig    ClusterPluginConfig `json:"new_config"`
	ReleasedSort int                 `json:"released_sort"`
	NowSort      int                 `json:"now_sort"`
	OptType      enum.ChangeOptType  `json:"opt_type"`
	CreateTime   string              `json:"create_time"`
}

// ClusterPluginToPublishItem 待发布集群插件列表项
type ClusterPluginToPublishItem struct {
	Name             string              `json:"name"`
	Extended         string              `json:"extended"`
	ReleasedConfig   ClusterPluginConfig `json:"released_config"`
	NoReleasedConfig ClusterPluginConfig `json:"no_released_config"`
	ReleasedSort     int                 `json:"released_sort"`
	NowSort          int                 `json:"now_sort"`
	CreateTime       string              `json:"create_time"`
	OptType          enum.ChangeOptType  `json:"opt_type"`
}

type ClusterPluginPublishInput struct {
	VersionName string `json:"version_name"`
	Desc        string `json:"desc"`
	Source      string `json:"source"`
}
