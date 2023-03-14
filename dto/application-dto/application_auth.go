package application_dto

import (
	"github.com/eolinker/apinto-dashboard/dto/discover-dto"
)

type ApplicationAuthInput struct {
	Driver        string                   `json:"driver"`
	IsTransparent bool                     `json:"is_transparent"`
	TokenName     string                   `json:"token_name"`
	Position      string                   `json:"position"`
	ExpireTime    int64                    `json:"expire_time"`
	Config        discover_dto.ConfigProxy `json:"config"`
}

type ApplicationAuthListOut struct {
	Uuid          string `json:"uuid"`
	Driver        string `json:"driver"`
	ParamPosition string `json:"param_position"`
	ParamName     string `json:"param_name"`
	ParamInfo     string `json:"param_info"`
	ExpireTime    int64  `json:"expire_time"`
	Operator      string `json:"operator"`
	UpdateTime    string `json:"update_time"`
	RuleInfo      string `json:"rule_info"`
	IsTransparent bool   `json:"is_transparent"`
}

type ApplicationAuthOut struct {
	Uuid          string          `json:"uuid"`
	Driver        string          `json:"driver"`
	ExpireTime    int64           `json:"expire_time"`
	Operator      string          `json:"operator"`
	Position      string          `json:"position"`
	TokenName     string          `json:"token_name"`
	UpdateTime    string          `json:"update_time"`
	IsTransparent bool            `json:"is_transparent"`
	Config        AuthConfigProxy `json:"config"`
}

type AuthConfigProxy []byte

func (c *AuthConfigProxy) MarshalJSON() ([]byte, error) {
	return *c, nil
}

func (c *AuthConfigProxy) UnmarshalJSON(bytes []byte) error {
	*c = bytes
	return nil
}
