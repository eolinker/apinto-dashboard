package dynamic_entry

import (
	"time"
)

type DynamicPublishVersion struct {
	Id          int                   `json:"id"`
	ClusterId   int                   `json:"cluster_id"`
	NamespaceId int                   `json:"namespace_id"`
	Publish     *DynamicPublishConfig `json:"publish"`
	Operator    int                   `json:"operator"`
	CreateTime  time.Time             `json:"create_time"`
}

func (v *DynamicPublishVersion) SetVersionId(id int) {
	v.Id = id
}

type DynamicPublishConfig struct {
	*BasicInfo
	Append map[string]interface{}
}

type BasicInfo struct {
	Profession  string `json:"profession"`
	Id          string `json:"id"`
	Name        string `json:"name"`
	Driver      string `json:"driver"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Create      string `json:"create"`
	Update      string `json:"update"`
}
