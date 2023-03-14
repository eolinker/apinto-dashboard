package application_model

import (
	"github.com/eolinker/apinto-dashboard/entry/application-entry"
)

const (
	AuthDriverBasic  = "basic"
	AuthDriverApikey = "apikey"
	AuthDriverAkSk   = "aksk"
	AuthDriverJwt    = "jwt"
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
