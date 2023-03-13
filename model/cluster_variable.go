package model

import (
	"github.com/eolinker/apinto-dashboard/entry"
	"time"
)

type ClusterVariable struct {
	*entry.ClusterVariable
}

type ClustersVariables struct {
	Clusters  []*entry.Cluster
	Variables []*entry.ClusterVariable
}

type ClusterVariableListItem struct {
	*entry.ClusterVariable
	Desc     string
	Operator string
	Publish  int //1.未发布 2.已发布 3.缺失
}

type ClusterVariableHistory struct {
	*entry.VariableHistory
}

type VariableToPublish struct {
	entry.VariableToPublish
}

type VariablePublishVersion entry.VariablePublishVersion

type VariablePublish struct {
	Id         int
	Name       string
	OptType    int //1.发布 2.回滚
	Operator   string
	CreateTime time.Time
	Details    []*VariablePublishDetails
}

type VariablePublishDetails struct {
	Key        string
	OldValue   string
	NewValue   string
	OptType    int //1.新增 2.修改 3.删除
	CreateTime time.Time
}

type ClusterVariableDiff struct {
	Key        string
	Value      string
	CreateTime time.Time
}

func (c *ClusterVariableDiff) GetKey() string {
	return c.Key
}

func (c *ClusterVariableDiff) Values() []string {
	return []string{c.Value}
}
