package online_dto

import "github.com/eolinker/apinto-dashboard/enum"

type OnlineOut struct {
	Name       string            `json:"name"`
	Status     enum.OnlineStatus `json:"status"`
	Env        string            `json:"env"`
	Disable    bool              `json:"disable"`
	Operator   string            `json:"operator"`
	UpdateTime string            `json:"update_time"`
}

type UpdateOnlineStatusInput struct {
	ClusterNames []string `json:"cluster_names"`
}
