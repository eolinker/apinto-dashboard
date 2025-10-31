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

type AkSkConfig struct {
	Ak    string            `json:"ak"`
	Sk    string            `json:"sk"`
	Label map[string]string `json:"label"`
}

type AkSk struct {
	apintoDriverName string
}

func (a *AkSk) GetAuthListInfo(config []byte) string {
	authConfig := &AuthConfig{}
	_ = json.Unmarshal(config, authConfig)
	return fmt.Sprintf("AK=%s,SK=%s", authConfig.Ak, authConfig.Sk)
}

func (a *AkSk) CheckInput(config []byte) ([]byte, error) {
	akSkConfig := new(AkSkConfig)

	if err := json.Unmarshal(config, akSkConfig); err != nil {
		return nil, err
	}
	if akSkConfig.Sk == "" || akSkConfig.Ak == "" {
		return nil, errors.New("ak or sk is null")
	}

	for key, _ := range akSkConfig.Label {
		if _, ok := enum.Keyword[strings.ToLower(key)]; ok {
			return nil, errors.New(fmt.Sprintf("标签信息中的%s为系统关键字", key))
		}
	}
	return config, nil
}

func (a *AkSk) GetCfgDetails(config []byte) []application_model.AuthDetailItem {
	akSkConfig := new(AkSkConfig)
	if err := json.Unmarshal(config, akSkConfig); err != nil {
		return []application_model.AuthDetailItem{}
	}
	return []application_model.AuthDetailItem{{Key: "AK", Value: akSkConfig.Ak}, {Key: "SK", Value: akSkConfig.Sk}}
}

func (a *AkSk) ToApinto(expire int64, position string, tokenName string, config []byte, hideCredential bool) v1.ApplicationAuth {

	akSkConfig := new(AkSkConfig)

	_ = json.Unmarshal(config, akSkConfig)

	akSkApintoConfig := &v1.ApplicationAuthAkSkConfig{
		Ak:     akSkConfig.Ak,
		Sk:     akSkConfig.Sk,
		Expire: expire,
	}

	users := make([]v1.ApplicationAuthUser, 0)
	pattern := &v1.ApplicationAuthUserAkSkPattern{
		Ak: akSkConfig.Ak,
		Sk: akSkConfig.Sk,
	}
	users = append(users, v1.ApplicationAuthUser{
		Expire:         expire,
		Labels:         akSkConfig.Label,
		Pattern:        pattern,
		HideCredential: hideCredential,
	})
	auth := v1.ApplicationAuth{
		Config:    akSkApintoConfig,
		Type:      a.apintoDriverName,
		Users:     users,
		Position:  position,
		TokenName: tokenName,
	}

	return auth
}

func CreateAkSk() application.IAuthDriver {
	return &AkSk{apintoDriverName: "aksk"}
}

func (a *AkSk) Render() string {
	return akSkConfigRender
}

var akSkConfigRender = `{
                "type": "object",
                "properties": {
                        "ak": {
                                "type": "string",
                                "title": "AK",
                                "x-component": "Input",
                                "x-component-props": {
                                        "placeholder": "请输入"
                                },
                                "required": true
                        },
                        "sk": {
                                "type": "string",
                                "title": "SK",
                                "x-component": "Input",
                                "x-component-props": {
                                        "placeholder": "请输入"
                                },
                                "required": true
                        }
                }
        }`
