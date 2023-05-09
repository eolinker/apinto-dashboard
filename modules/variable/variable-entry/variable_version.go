package variable_entry

import "time"

// VariablePublishVersion 集群绑定的环境变量版本信息表
type VariablePublishVersion struct {
	Id          int `json:"id"`
	ClusterId   int `json:"cluster_id"`
	NamespaceId int `json:"namespace_id"`
	VariablePublishVersionConfig
	Desc       string    `json:"desc"`
	Operator   int       `json:"operator"`
	CreateTime time.Time `json:"create_time"`
}

func (v *VariablePublishVersion) SetVersionId(id int) {
	v.Id = id
}

type VariablePublishVersionConfig struct {
	ClusterVariable []*ClusterVariable `json:"cluster_variable"`
}

type VariableToPublish struct {
	Key             string    `json:"key"`
	VariableId      int       `json:"variable_id"`
	FinishValue     string    `json:"finish_value"`
	NoReleasedValue string    `json:"no_released_value"`
	CreateTime      time.Time `json:"create_time"`
	OptType         int       `json:"opt_type"` //操作类型(1新增 2更新 3删除)
}
