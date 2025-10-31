package driver

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/eolinker/eosc/log"

	"golang.org/x/crypto/pbkdf2"

	v1 "github.com/eolinker/apinto-dashboard/client/v1"
	"github.com/eolinker/apinto-dashboard/modules/application"
	application_model "github.com/eolinker/apinto-dashboard/modules/application/application-model"
)

type Oauth2 struct {
	driver string
}

type Oauth2Config struct {
	*v1.ApplicationAuthUserOauth2Pattern
	Hashed bool `json:"hashed"`
}

func (j *Oauth2) CheckInput(config []byte) ([]byte, error) {
	cfg := new(Oauth2Config)
	if err := json.Unmarshal(config, cfg); err != nil {
		return nil, err
	}
	if cfg.HashSecret && !cfg.Hashed {
		// 未加密
		secret, err := hashSecret([]byte(cfg.ClientSecret), 0, 0, 0)
		if err != nil {
			log.Error(err)
		} else {
			cfg.ClientSecret = secret
			cfg.Hashed = true
		}
	}
	return json.Marshal(cfg)
}

func (j *Oauth2) GetCfgDetails(config []byte) []application_model.AuthDetailItem {
	cfg := new(Oauth2Config)
	if err := json.Unmarshal(config, cfg); err != nil {
		return []application_model.AuthDetailItem{}
	}
	redirectURLs, _ := json.Marshal(cfg.RedirectUrls)
	items := []application_model.AuthDetailItem{
		{Key: "客户端ID", Value: cfg.ClientId},
		{Key: "客户端密钥", Value: cfg.ClientSecret},
		{Key: "客户端类型", Value: cfg.ClientType},
		{Key: "对密钥进行Hash", Value: strconv.FormatBool(cfg.HashSecret)},
		{Key: "重定向URL列表", Value: string(redirectURLs)},
	}

	return items
}

func (j *Oauth2) GetAuthListInfo(config []byte) string {
	authConfig := &AuthConfig{}
	_ = json.Unmarshal(config, authConfig)
	return authConfig.User
}

func (j *Oauth2) ToApinto(expire int64, position string, tokenName string, config []byte, hideCredential bool) v1.ApplicationAuth {

	cfg := new(Oauth2Config)

	_ = json.Unmarshal(config, cfg)

	users := make([]v1.ApplicationAuthUser, 0)
	users = append(users, v1.ApplicationAuthUser{
		Expire:         expire,
		Pattern:        cfg.ApplicationAuthUserOauth2Pattern,
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

func CreateOauth2() application.IAuthDriver {
	return &Oauth2{driver: "oauth2"}
}

func (j *Oauth2) Render() string {
	return oauth2ConfigRender
}

func hashSecret(secret []byte, saltLen int, iterations int, keyLength int) (string, error) {
	if saltLen < 1 {
		saltLen = 16
	}
	salt, err := generateRandomSalt(saltLen)
	if err != nil {
		return "", err
	}
	// 迭代次数和密钥长度
	if iterations < 1 {
		iterations = 10000
	}
	if keyLength < 1 {
		keyLength = 32
	}

	// 使用 PBKDF2 密钥派生函数
	key := pbkdf2.Key(secret, salt, iterations, keyLength, sha512.New)
	return fmt.Sprintf("$pbkdf2-sha512$i=%d,l=%d$%s$%s", iterations, keyLength, base64.RawStdEncoding.EncodeToString(salt), base64.RawStdEncoding.EncodeToString(key)), nil
}

func generateRandomSalt(length int) ([]byte, error) {
	// Create a byte slice with the specified length
	salt := make([]byte, length)

	// Use crypto/rand to fill the slice with random bytes
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}

	// Return the salt as a hexadecimal string
	return salt, nil
}

var oauth2ConfigRender = `{
    "type": "object",
    "properties": {
        "client_id": {
            "type": "string",
            "title": "客户端ID",
            "x-component": "Input",
            "x-component-props": {
                "placeholder": "请输入"
            },
            "required": true
        },
        "client_secret": {
            "type": "string",
            "title": "客户端密钥",
            "x-component": "Input",
            "x-component-props": {
                "placeholder": "请输入"
            },
            "required": true
        },
        "hash_secret": {
            "type": "boolean",
            "title": "对密钥进行Hash",
            "x-component": "Checkbox",
            "label": "",
            "default": true
        },
        "redirect_urls": {
            "type": "array",
            "title": "重定向URL列表",
            "x-component": "ArrayItems",
            "items": {
                "type": "void",
                "x-component": "Space",
                "x-component-props": {
                    "placeholder": "请输入"
                },
                "properties": {
                    "url": {
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
                "redirectUris0": {
                    "type": "void",
                    "x-component": "Space",
                    "x-component-props": {
                        "placeholder": "请输入"
                    },
                    "x-index": 0,
                    "properties": {
                        "url": {
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
        },
        "client_type": {
            "enum": [
                {
                    "label": "confidential",
                    "value": "confidential"
                },
                {
                    "label": "public",
                    "value": "public"
                }
            ],
            "type": "string",
            "default": "confidential",
            "title": "客户端类型",
            "required": true,
            "x-component": "Select"
        }
    }
}`
