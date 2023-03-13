package model

import (
	"github.com/eolinker/apinto-dashboard/entry"
)

const (
	AuthDriverBasic  = "basic"
	AuthDriverApikey = "apikey"
	AuthDriverAkSk   = "aksk"
	AuthDriverJwt    = "jwt"
)

type ApplicationAuth struct {
	*entry.ApplicationAuth
	Operator      string
	ParamPosition string
	ParamName     string
	ParamInfo     string
	RuleInfo      string //规则信息
	Config        string //配置信息
}
