package driver

import (
	"encoding/json"
	"errors"
	"fmt"
	v1 "github.com/eolinker/apinto-dashboard/client/v1"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/enum"
	"strings"
)

type Basic struct {
	apintoDriverName string
}

type BasicConfig struct {
	UserName string            `json:"user_name"`
	Password string            `json:"password"`
	Label    map[string]string `json:"label"`
}

func (b *Basic) CheckInput(config []byte) error {
	basicConfig := new(BasicConfig)
	if err := json.Unmarshal(config, basicConfig); err != nil {
		return err
	}
	if basicConfig.UserName == "" || basicConfig.Password == "" {
		return errors.New("username or password is null")
	}

	for key, _ := range basicConfig.Label {
		if _, ok := enum.Keyword[strings.ToLower(key)]; ok {
			return errors.New(fmt.Sprintf("标签信息中的%s为系统关键字", key))
		}
	}

	return nil
}

func (b *Basic) GetAuthListInfo(config []byte) string {
	authConfig := &AuthConfig{}
	_ = json.Unmarshal(config, authConfig)
	secret := fmt.Sprintf("%s:%s", authConfig.UserName, authConfig.Password)
	return common.Base64Encode([]byte(secret))
}

func (b *Basic) ToApinto(expire int64, position string, tokenName string, config []byte, hideCredential bool) v1.ApplicationAuth {

	basicConfig := new(BasicConfig)

	_ = json.Unmarshal(config, basicConfig)

	basicApintoConfig := &v1.ApplicationAuthBasicConfig{
		Username: basicConfig.UserName,
		Password: basicConfig.Password,
		Expire:   expire,
	}

	users := make([]v1.ApplicationAuthUser, 0)
	pattern := &v1.ApplicationAuthUserBasicPattern{
		Username: basicConfig.UserName,
		Password: basicConfig.Password,
	}
	users = append(users, v1.ApplicationAuthUser{
		Expire:         expire,
		Labels:         basicConfig.Label,
		Pattern:        pattern,
		HideCredential: hideCredential,
	})
	auth := v1.ApplicationAuth{
		Config:    basicApintoConfig,
		Type:      b.apintoDriverName,
		Users:     users,
		Position:  position,
		TokenName: tokenName,
	}

	return auth
}

func CreateBasic() IAuthDriver {
	return &Basic{apintoDriverName: "basic"}
}

func (b *Basic) Render() string {
	return basicConfigRender
}

var basicConfigRender = `{
	"type": "object",
	"properties": {
		"user_name": {
			"type": "string",
			"title": "用户名",
			"x-component": "Input",
			"x-component-props": {
				"placeholder": "请输入",
				"extra": "英文数字下划线任意一种，首字母必须为英文"
			},
			"pattern": "^[a-zA-Z][a-zA-Z0-9]*",
			"required": true
		},
		"password": {
			"type": "string",
			"title": "密码",
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
				"x-component-props": {
					"placeholder": "请输入"
				},
				"properties": {
					"key": {
						"type": "string",
						"x-component": "Input",
						"x-component-props": {
							"class": "w240",
							"placeholder": "请输入key"
						},
						"x-index": 0
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
