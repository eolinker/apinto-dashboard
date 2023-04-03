package open_app_model

import (
	"github.com/eolinker/apinto-dashboard/modules/openapp/open-app-entry"
)

type ExternalAppInfo struct {
	*open_app_entry.ExternalApplication
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
