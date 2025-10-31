package driver

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	v1 "github.com/eolinker/apinto-dashboard/client/v1"
	"github.com/eolinker/apinto-dashboard/modules/application"
	application_model "github.com/eolinker/apinto-dashboard/modules/application/application-model"
)

type Jwt struct {
	apintoDriverName string
}

type JwtConfig struct {
	Iss               string            `json:"iss"`
	Algorithm         string            `json:"algorithm"`
	Secret            string            `json:"secret"`
	PublicKey         string            `json:"public_key"`
	User              string            `json:"user"`
	UserPath          string            `json:"user_path"`
	ClaimsToVerify    []string          `json:"claims_to_verify"`
	SignatureIsBase64 bool              `json:"signature_is_base64"`
	Label             map[string]string `json:"label"`
}

func (j *Jwt) CheckInput(config []byte) ([]byte, error) {

	jwtConfig := new(JwtConfig)
	if err := json.Unmarshal(config, jwtConfig); err != nil {
		return nil, err
	}
	if jwtConfig.Iss == "" {
		return nil, errors.New("iss is null")
	}
	if jwtConfig.Algorithm == "" {
		return nil, errors.New("algorithm is null")
	}
	algorithm := strings.ToUpper(jwtConfig.Algorithm)
	switch algorithm {
	case "HS256", "HS384", "HS512":
		if jwtConfig.Secret == "" {
			return nil, errors.New("secret is null")
		}
	default:
		if jwtConfig.PublicKey == "" {
			return nil, errors.New("public_key is null")
		}
	}

	//校验 校验字段
	for _, claim := range jwtConfig.ClaimsToVerify {
		switch claim {
		case "exp", "nbf":
		default:
			return nil, fmt.Errorf("claim key %s is illegal. ", claim)
		}
	}

	return config, nil
}

func (j *Jwt) GetCfgDetails(config []byte) []application_model.AuthDetailItem {
	jwtConfig := new(JwtConfig)
	if err := json.Unmarshal(config, jwtConfig); err != nil {
		return []application_model.AuthDetailItem{}
	}
	items := []application_model.AuthDetailItem{
		{Key: "Iss", Value: jwtConfig.Iss},
		{Key: "签名算法", Value: jwtConfig.Algorithm},
		{Key: "用户名", Value: jwtConfig.User},
		{Key: "用户名JsonPath", Value: jwtConfig.UserPath},
		{Key: "校验字段", Value: strings.Join(jwtConfig.ClaimsToVerify, ",")},
	}

	switch jwtConfig.Algorithm {
	case "HS256", "HS384", "HS512":
		items = append(items, application_model.AuthDetailItem{Key: "Secret", Value: jwtConfig.Secret})
		base64 := "false"
		if jwtConfig.SignatureIsBase64 {
			base64 = "true"
		}
		items = append(items, application_model.AuthDetailItem{Key: "Secret", Value: base64})
	default:
		items = append(items, application_model.AuthDetailItem{Key: "RSA公钥", Value: jwtConfig.PublicKey})
	}

	return items
}

func (j *Jwt) GetAuthListInfo(config []byte) string {
	authConfig := &AuthConfig{}
	_ = json.Unmarshal(config, authConfig)
	return authConfig.User
}

func (j *Jwt) ToApinto(expire int64, position string, tokenName string, config []byte, hideCredential bool) v1.ApplicationAuth {

	jwtConfig := new(JwtConfig)

	_ = json.Unmarshal(config, jwtConfig)

	basicApintoConfig := &v1.ApplicationAuthJwtConfig{
		Iss:               jwtConfig.Iss,
		Secret:            jwtConfig.Secret,
		RsaPublicKey:      jwtConfig.PublicKey,
		Algorithm:         jwtConfig.Algorithm,
		ClaimsToVerify:    jwtConfig.ClaimsToVerify,
		SignatureIsBase64: jwtConfig.SignatureIsBase64,
		Path:              jwtConfig.UserPath,
	}

	users := make([]v1.ApplicationAuthUser, 0)
	pattern := &v1.ApplicationAuthUserJwtPattern{
		UserName: jwtConfig.User,
	}
	users = append(users, v1.ApplicationAuthUser{
		Expire:         expire,
		Labels:         jwtConfig.Label,
		Pattern:        pattern,
		HideCredential: hideCredential,
	})
	auth := v1.ApplicationAuth{
		Config:    basicApintoConfig,
		Type:      j.apintoDriverName,
		Users:     users,
		Position:  position,
		TokenName: tokenName,
	}

	return auth
}

func CreateJwt() application.IAuthDriver {
	return &Jwt{apintoDriverName: "jwt"}
}

func (j *Jwt) Render() string {
	return jwtConfigRender
}

var jwtConfigRender = `{
	"type": "object",
	"properties": {
		"algorithm": {
			"enum": [{
					"label": "HS256",
					"value": "HS256"
				},
				{
					"label": "HS384",
					"value": "HS384"
				},
				{
					"label": "HS512",
					"value": "HS512"
				},
				{
					"label": "RS256",
					"value": "RS256"
				},
				{
					"label": "RS384",
					"value": "RS384"
				},
				{
					"label": "RS512",
					"value": "RS512"
				},
				{
					"label": "ES256",
					"value": "ES256"
				},
				{
					"label": "ES384",
					"value": "ES384"
				},
				{
					"label": "ES512",
					"value": "ES512"
				}
			],
			"type": "string",
			"default": "HS256",
			"title": "算法",
			"required": true,
			"x-component": "Select"
		},
		"iss": {
			"type": "string",
			"title": "Iss",
			"x-component": "Input",
			"x-component-props": {
				"placeholder": "请输入"
			},
			"required": true
		},
		"secret": {
			"type": "string",
			"title": "Secret",
			"x-component": "Input",
			"x-component-props": {
				"placeholder": "请输入"
			},
			"required": true,
			"x-reactions": {
				"dependencies": [
					"algorithm"
				],
				"fulfill": {
					"state": {
						"display": "{{$deps[0].include('HS')}}"
					}
				}
			}
		},
		"public_key": {
			"type": "string",
			"title": "RsaPublicKey",
			"x-component": "Input",
			"x-component-props": {
				"placeholder": "请输入"
			},
			"required": true,
			"x-reactions": {
				"dependencies": [
					"algorithm"
				],
				"fulfill": {
					"state": {
						"display": "{{$deps[0].include('RS') || $deps[0].include('ES')}}"
					}
				}
			}
		},
		"label": {
			"type": "array",
			"title": "标签信息",
			"x-component": "ArrayItems",
			"items": {
				"type": "void",
				"x-component": "Space",
				"x-component-props": {
					"placeholder": "请输入"
				},
				"properties": {
					"key": {
						"type": "string",
						"x-component": "Input",
						"x-index": 0,
						"x-component-props": {
							"class": "w240",
							"placeholder": "请输入key"
						}
					},
					"value": {
						"type": "text",
						"x-component": "Input",
						"x-index": 1,
						"x-component-props": {
							"class": "w240 mg_button",
							"placeholder": "请输入value"
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
				"params0": {
					"type": "void",
					"x-component": "Space",
					"x-component-props": {
						"placeholder": "请输入"
					},
					"x-index": 0,
					"properties": {
						"key": {
							"type": "text",
							"x-component": "Input",
							"x-index": 0,
							"x-component-props": {
								"class": "w240",
								"placeholder": "请输入key"
							}
						},
						"value": {
							"type": "text",
							"x-component": "Input",
							"x-index": 1,
							"x-component-props": {
								"class": "w240 mg_button",
								"placeholder": "请输入value"
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
				},
				"params1": {
					"type": "void",
					"x-component": "Space",
					"x-component-props": {
						"placeholder": "请输入"
					},
					"x-index": 0,
					"properties": {
						"key": {
							"type": "text",
							"x-component": "Input",
							"x-index": 0,
							"x-component-props": {
								"class": "w240",
								"placeholder": "请输入key"
							}
						},
						"value": {
							"type": "text",
							"x-component": "Input",
							"x-index": 1,
							"x-component-props": {
								"class": "w240 mg_button",
								"placeholder": "请输入value"
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
							"x-component-props": {
								"class": "mg_button"
							},
							"x-index": 2
						}
					}
				}
			}
		}
	}
}`
