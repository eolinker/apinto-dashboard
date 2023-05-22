package driver

import (
	"encoding/json"
	"errors"
	"fmt"
	v1 "github.com/eolinker/apinto-dashboard/client/v1"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/modules/application"
	"strings"
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

func (a *AkSk) CheckInput(config []byte) error {

	akSkConfig := new(AkSkConfig)

	if err := json.Unmarshal(config, akSkConfig); err != nil {
		return err
	}
	if akSkConfig.Sk == "" || akSkConfig.Ak == "" {
		return errors.New("ak or sk is null")
	}

	for key, _ := range akSkConfig.Label {
		if _, ok := enum.Keyword[strings.ToLower(key)]; ok {
			return errors.New(fmt.Sprintf("标签信息中的%s为系统关键字", key))
		}
	}

	return nil
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
		},
		"label": {
			"type": "array",
			"title": "标签信息",
			"x-component": "ArrayItems",
			"items": {
				"type": "void",
				"x-component": "Space",
				"x-component-props": {},
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
