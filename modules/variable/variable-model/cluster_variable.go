package variable_model

import (
	"github.com/eolinker/apinto-dashboard/modules/cluster/cluster-entry"
	variable_entry2 "github.com/eolinker/apinto-dashboard/modules/variable/variable-entry"
	"time"
)

type ClusterVariable struct {
	*variable_entry2.ClusterVariable
}

type ClustersVariables struct {
	Clusters  []*cluster_entry.Cluster
	Variables []*variable_entry2.ClusterVariable
}

type ClusterVariableListItem struct {
	*variable_entry2.ClusterVariable
	Desc     string
	Operator string
	Publish  int //1.未发布 2.已发布 3.缺失
}

type ClusterVariableHistory struct {
	*variable_entry2.VariableHistory
}

type VariableToPublish struct {
	variable_entry2.VariableToPublish
}

type VariablePublishVersion variable_entry2.VariablePublishVersion

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
