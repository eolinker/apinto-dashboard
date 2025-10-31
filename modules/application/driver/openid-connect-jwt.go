package driver

import (
	"encoding/json"

	v1 "github.com/eolinker/apinto-dashboard/client/v1"
	"github.com/eolinker/apinto-dashboard/modules/application"
	application_model "github.com/eolinker/apinto-dashboard/modules/application/application-model"
)

type OpenIDConnectJWT struct {
	driver string
}

type OpenIDConnectJWTConfig struct {
	*v1.ApplicationAuthOpenIDConnectConfig
}

func (j *OpenIDConnectJWT) CheckInput(config []byte) ([]byte, error) {
	cfg := new(OpenIDConnectJWTConfig)
	if err := json.Unmarshal(config, cfg); err != nil {
		return nil, err
	}

	return config, nil
}

func (j *OpenIDConnectJWT) GetCfgDetails(config []byte) []application_model.AuthDetailItem {
	cfg := new(OpenIDConnectJWTConfig)
	if err := json.Unmarshal(config, cfg); err != nil {
		return []application_model.AuthDetailItem{}
	}
	claim, _ := json.Marshal(cfg.AuthenticatedGroupsClaim)
	items := []application_model.AuthDetailItem{
		{Key: "Issuer", Value: cfg.Issuer},
		{Key: "用户组提取字段", Value: string(claim)},
	}

	return items
}

func (j *OpenIDConnectJWT) GetAuthListInfo(config []byte) string {
	authConfig := &AuthConfig{}
	_ = json.Unmarshal(config, authConfig)
	return authConfig.User
}

func (j *OpenIDConnectJWT) ToApinto(expire int64, position string, tokenName string, config []byte, hideCredential bool) v1.ApplicationAuth {

	cfg := new(OpenIDConnectJWTConfig)

	_ = json.Unmarshal(config, cfg)

	users := make([]v1.ApplicationAuthUser, 0)
	users = append(users, v1.ApplicationAuthUser{
		Expire:         expire,
		Pattern:        cfg.ApplicationAuthOpenIDConnectConfig,
		HideCredential: hideCredential,
	})
	auth := v1.ApplicationAuth{
		Type:      j.driver,
		Users:     users,
		Position:  position,
		TokenName: tokenName,
	}

	return auth
}

func CreateOpenidConnectJWT() application.IAuthDriver {
	return &OpenIDConnectJWT{driver: "openid-connect-jwt"}
}

func (j *OpenIDConnectJWT) Render() string {
	return openidConnectJWTConfigRender
}

var openidConnectJWTConfigRender = `{
  "type": "object",
  "properties": {
    "issuer": {
      "type": "string",
      "title": "Issuer",
      "x-component": "Input",
      "x-component-props": {
        "placeholder": "请输入"
      },
      "required": true
    },
    "authenticated_groups_claim": {
      "type": "array",
      "title": "用户组提取字段",
      "required": true,
      "x-component": "ArrayItems",
      "items": {
        "type": "void",
        "x-component": "Space",
        "x-component-props": {
          "placeholder": "请输入"
        },
        "properties": {
          "value": {
            "type": "text",
            "x-component": "Input",
            "x-index": 1,
            "x-component-props": {
              "class": "w240 mg_button",
              "placeholder": "请输入"
            }
          },
          "remove": {
            "type": "void",
            "x-component": "ArrayItems.Remove",
            "x-index": 3
          },
          "add": {
            "type": "void",
            "x-component": "ArrayItems.Addition",
            "x-index": 2,
            "x-component-props": {
              "class": "mg_button"
            }
          }
        }
      },
      "properties": {
        "authenticatedGroupsClaim0": {
          "type": "void",
          "x-component": "Space",
          "x-component-props": {
            "placeholder": "请输入"
          },
          "x-index": 0,
          "properties": {
            "value": {
              "type": "text",
              "x-component": "Input",
              "x-index": 1,
              "x-component-props": {
                "class": "w240 mg_button",
                "placeholder": "请输入"
              }
            },
            "add": {
              "type": "void",
              "x-component": "ArrayItems.Addition",
              "x-component-props": {
                "class": "mg_button"
              },
              "x-index": 3
            }
          }
        }
      }
    }
  }
}`
