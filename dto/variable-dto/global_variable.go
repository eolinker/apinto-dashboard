package variable_dto

import "github.com/eolinker/apinto-dashboard/enum"

type GlobalVariableListItem struct {
	Key         string                   `json:"key"`
	Status      enum.VariableUsageStatus `json:"status"`
	Description string                   `json:"description"`
	Operator    string                   `json:"operator"`
	CreateTime  string                   `json:"create_time"`
}

type GlobalVariableDetailsItem struct {
	ClusterName string                      `json:"cluster_name"`
	Environment string                      `json:"environment"`
	Value       string                      `json:"value"`
	Status      enum.ClusterVariablePublish `json:"publish_status"`
}

type GlobalVariableInput struct {
	Key  string `json:"key"`
	Desc string `json:"desc"`
}
