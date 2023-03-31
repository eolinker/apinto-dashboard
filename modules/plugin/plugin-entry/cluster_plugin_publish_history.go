package plugin_entry

import "time"

type ClusterPluginPublishHistory struct {
	Id          int
	VersionName string
	ClusterId   int
	NamespaceId int
	Desc        string
	VersionId   int
	CluPluginPublishHistoryInfo
	OptType  int
	Operator int
	OptTime  time.Time
}

type CluPluginPublishHistoryInfo struct {
	PluginToPublish []*PluginToPublish `json:"plugin_to_publish"` //发布记录
}

type PluginToPublish struct {
	PluginName       string                `json:"plugin_name"`
	Extended         string                `json:"extended"`
	Rely             string                `json:"rely"`
	ReleasedConfig   PluginToPublishConfig `json:"released_config"`
	NoReleasedConfig PluginToPublishConfig `json:"no_released_config"`
	ReleasedSort     int                   `json:"released_sort"`
	NowSort          int                   `json:"now_sort"`
	CreateTime       time.Time             `json:"create_time"`
	OptType          int                   `json:"opt_type"` //操作类型(1新增 2更新 3删除)
}

type PluginToPublishConfig struct {
	Status int    `json:"status"`
	Config string `json:"config"`
}
