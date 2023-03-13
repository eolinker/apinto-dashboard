package dto

import "github.com/eolinker/apinto-dashboard/enum"

type ClusterVariableItem struct {
	Id         int                         `json:"id,omitempty"`          //环境变量主键ID
	VariableId int                         `json:"variable_id,omitempty"` //全局环境变量ID
	Key        string                      `json:"key"`
	Value      string                      `json:"value"`
	Publish    enum.ClusterVariablePublish `json:"publish"`
	Desc       string                      `json:"desc"`
	Operator   string                      `json:"operator"`
	UpdateTime string                      `json:"update_time"`
}

type ClusterVariableSyncConf struct {
	Id         int    `json:"id,omitempty"`          //环境变量主键ID
	VariableId int    `json:"variable_id,omitempty"` //全局环境变量ID
	Key        string `json:"key"`
	Value      string `json:"value"`
	Desc       string `json:"desc"`
	UpdateTime string `json:"update_time"`
}

type SyncConf struct {
	Clusters  []*ClusterInput            `json:"clusters"`
	Variables []*ClusterVariableSyncConf `json:"variables"`
}

type ClusterHistoryOut struct {
	Key        string               `json:"key"`
	OldValue   string               `json:"old_value"`
	NewValue   string               `json:"new_value"`
	CreateTime string               `json:"create_time"`
	OptType    enum.VariableOptType `json:"opt_type"`
}

type VariableToPublishOut struct {
	Key             string               `json:"key"`
	FinishValue     string               `json:"finish_value"`
	NoReleasedValue string               `json:"no_released_value"`
	CreateTime      string               `json:"create_time"`
	OptType         enum.VariableOptType `json:"opt_type"`
}

type VariablePublishInput struct {
	VersionName string `json:"version_name"`
	Desc        string `json:"desc"`
	Source      string `json:"source"`
}

type VariablePublishOut struct {
	Id         int                       `json:"id"`
	Name       string                    `json:"name"`
	OptType    enum.PublishOptType       `json:"opt_type"`
	Operator   string                    `json:"operator"`
	CreateTime string                    `json:"create_time"`
	Details    []*VariablePublishDetails `json:"details"`
}

type VariablePublishDetails struct {
	Key        string               `json:"key"`
	OldValue   string               `json:"old_value"`
	NewValue   string               `json:"new_value"`
	OptType    enum.VariableOptType `json:"opt_type"`
	CreateTime string               `json:"create_time"`
}
