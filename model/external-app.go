package model

import "github.com/eolinker/apinto-dashboard/entry"

type ExternalAppInfo struct {
	*entry.ExternalApplication
}

type ExtAppListItem struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Token      string `json:"token"`
	Tags       string `json:"tags"`
	Status     int    `json:"status"`
	Operator   string `json:"operator"`
	UpdateTime string `json:"update_time"`
}
