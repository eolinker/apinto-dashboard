package dto

import "github.com/eolinker/apinto-dashboard/enum"

type ServiceApiInput struct {
	Method    string `json:"method"`
	Uri       string `json:"uri"`
	GroupUUID string `json:"group_uuid"`
	Name      string `json:"name"`
	Config    string `json:"config"`
}

type DeleteServiceApiInput struct {
	Ids []string `json:"ids"`
}

type ServiceApiOut struct {
	Method     enum.ServiceApiMethod `json:"method"`
	Uri        string                `json:"uri"`
	Name       string                `json:"name"`
	Config     string                `json:"config"`
	Id         string                `json:"id"`
	Operator   string                `json:"operator"`
	UpdateTime string                `json:"update_time"`
	IsDelete   bool                  `json:"is_delete"`
}
