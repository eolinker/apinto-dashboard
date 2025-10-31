package application_model

import (
	"github.com/eolinker/apinto-dashboard/modules/application/application-entry"
	"time"
)

type ApplicationAuth struct {
	*application_entry.ApplicationAuth
	Operator      string
	ParamPosition string
	ParamName     string
	ParamInfo     string
	RuleInfo      string //规则信息
	Config        string //配置信息
}

func (a *ApplicationAuth) UserId() int {
	return a.ApplicationAuth.Operator
}

func (a *ApplicationAuth) Set(name string) {
	a.Operator = name
}

type AppAuthItem struct {
	UUID           string
	Title          string
	Driver         string
	Operator       string
	HideCredential bool
	ExpireTime     int64
	UpdateTime     time.Time
}

type AuthDetailItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
