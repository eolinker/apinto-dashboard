package driver

import (
	"encoding/json"
	"errors"
	"fmt"
	v1 "github.com/eolinker/apinto-dashboard/client/v1"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/modules/application"
	application_model "github.com/eolinker/apinto-dashboard/modules/application/application-model"
	"strings"
)

type ApikeyConfig struct {
	Apikey string            `json:"apikey"`
	Label  map[string]string `json:"label"`
}

type Apikey struct {
	apintoDriverName string
}

func (a *Apikey) GetAuthListInfo(config []byte) string {
	authConfig := &AuthConfig{}
	_ = json.Unmarshal(config, authConfig)
	return authConfig.Apikey
}

func (a *Apikey) CheckInput(config []byte) error {
	apikeyConfig := new(ApikeyConfig)
	if err := json.Unmarshal(config, apikeyConfig); err != nil {
		return err
	}
	if apikeyConfig.Apikey == "" {
		return errors.New("apikey is null")
	}

	for key, _ := range apikeyConfig.Label {
		if _, ok := enum.Keyword[strings.ToLower(key)]; ok {
			return errors.New(fmt.Sprintf("标签信息中的%s为系统关键字", key))
		}
	}

	return nil
}

func (a *Apikey) GetCfgDetails(config []byte) []application_model.AuthDetailItem {
	apiKeyConfig := new(ApikeyConfig)
	if err := json.Unmarshal(config, apiKeyConfig); err != nil {
		return []application_model.AuthDetailItem{}
	}
	return []application_model.AuthDetailItem{{Key: "Apikey", Value: apiKeyConfig.Apikey}}
}

func (a *Apikey) ToApinto(expire int64, position string, tokenName string, config []byte, hideCredential bool) v1.ApplicationAuth {

	apikeyConfig := new(ApikeyConfig)

	_ = json.Unmarshal(config, apikeyConfig)

	apikeyApintoConfig := &v1.ApplicationAuthApiKeyConfig{
		Apikey: apikeyConfig.Apikey,
		Expire: expire,
	}

	users := make([]v1.ApplicationAuthUser, 0)
	pattern := &v1.ApplicationAuthUserApiKeyPattern{
		Apikey: apikeyConfig.Apikey,
	}
	users = append(users, v1.ApplicationAuthUser{
		Expire:         expire,
		Labels:         apikeyConfig.Label,
		Pattern:        pattern,
		HideCredential: hideCredential,
	})
	auth := v1.ApplicationAuth{
		Config:    apikeyApintoConfig,
		Type:      a.apintoDriverName,
		Users:     users,
		Position:  position,
		TokenName: tokenName,
	}

	return auth
}

func CreateApikey() application.IAuthDriver {
	return &Apikey{apintoDriverName: "apikey"}
}

func (a *Apikey) Render() string {
	return apikeyConfigRender
}

var apikeyConfigRender = `{
                "type": "object",
                "properties": {
                        "apikey": {
                                "type": "string",
                                "title": "Apikey",
                                "x-component": "Input",
                                "x-component-props": {
                                        "placeholder": "请输入"
                                },
                                "required": true
                        }
                }
        }`
