package driver

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	v1 "github.com/eolinker/apinto-dashboard/client/v1"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/modules/application"
	application_model "github.com/eolinker/apinto-dashboard/modules/application/application-model"
)

type ParaHmacConfig struct {
	AppID  string            `json:"app_id"`
	AppKey string            `json:"app_key"`
	Label  map[string]string `json:"label"`
}

type ParaHmac struct {
	driver string
}

func (a *ParaHmac) GetAuthListInfo(config []byte) string {
	authConfig := &AuthConfig{}
	_ = json.Unmarshal(config, authConfig)
	return fmt.Sprintf("AppId=%s,AppKey=%s", authConfig.AppID, authConfig.AppKey)
}

func (a *ParaHmac) CheckInput(config []byte) ([]byte, error) {
	cfg := new(ParaHmacConfig)

	if err := json.Unmarshal(config, cfg); err != nil {
		return nil, err
	}
	if cfg.AppID == "" || cfg.AppKey == "" {
		return nil, errors.New("app_id or app_key is null")
	}

	for key := range cfg.Label {
		if _, ok := enum.Keyword[strings.ToLower(key)]; ok {
			return nil, errors.New(fmt.Sprintf("标签信息中的%s为系统关键字", key))
		}
	}
	return config, nil
}

func (a *ParaHmac) GetCfgDetails(config []byte) []application_model.AuthDetailItem {
	cfg := new(ParaHmacConfig)
	if err := json.Unmarshal(config, cfg); err != nil {
		return []application_model.AuthDetailItem{}
	}
	return []application_model.AuthDetailItem{{Key: "APP ID", Value: cfg.AppID}, {Key: "APP Key", Value: cfg.AppKey}}
}

func (a *ParaHmac) ToApinto(expire int64, position string, tokenName string, config []byte, hideCredential bool) v1.ApplicationAuth {

	cfg := new(ParaHmacConfig)

	_ = json.Unmarshal(config, cfg)

	pattern := &v1.ApplicationAuthUserParaHmacPattern{
		AppId:  cfg.AppID,
		AppKey: cfg.AppKey,
	}
	users := make([]v1.ApplicationAuthUser, 0)
	users = append(users, v1.ApplicationAuthUser{
		Expire:         expire,
		Labels:         cfg.Label,
		Pattern:        pattern,
		HideCredential: hideCredential,
	})
	auth := v1.ApplicationAuth{
		Config:    nil,
		Type:      a.driver,
		Users:     users,
		Position:  position,
		TokenName: tokenName,
	}

	return auth
}

func CreateParaHmac() application.IAuthDriver {
	return &ParaHmac{driver: "para-hmac"}
}

func (a *ParaHmac) Render() string {
	return paraHmacConfigRender
}

var paraHmacConfigRender = `{
                "type": "object",
                "properties": {
                        "app_id": {
                                "type": "string",
                                "title": "APP ID",
                                "x-component": "Input",
                                "x-component-props": {
                                        "placeholder": "请输入"
                                },
                                "required": true
                        },
                        "app_key": {
                                "type": "string",
                                "title": "APP Key",
                                "x-component": "Input",
                                "x-component-props": {
                                        "placeholder": "请输入"
                                },
                                "required": true
                        }
                }
        }`
