package v2

import (
	"encoding/json"
	"time"
)

type Config struct {
	Type      string            `json:"type"`
	Position  string            `json:"position"`
	Users     *User             `json:"users"`
	TokenName string            `json:"token_name"`
	Labels    map[string]string `json:"labels,omitempty"`
}

func (c *Config) MarshalJSON() ([]byte, error) {
	tmp := make(map[string]interface{})
	if c.Users != nil {
		c.Users.Labels = c.Labels
	}
	tmp["type"] = c.Type
	tmp["position"] = c.Position
	tmp["token_name"] = c.TokenName
	tmp["users"] = []*User{
		c.Users,
	}
	return json.Marshal(tmp)
}

type User struct {
	Expire         int64             `json:"expire"`
	Pattern        map[string]string `json:"pattern"`
	HideCredential bool              `json:"hide_credential"`
	Labels         map[string]string `json:"labels"`
}

func (u *User) UnmarshalJSON(bytes []byte) error {
	var tmp map[string]interface{}
	err := json.Unmarshal(bytes, &tmp)
	if err != nil {
		return err
	}
	if vv, has := tmp["expire"]; has {
		v, ok := vv.(string)
		if ok {
			t, err := time.ParseInLocation("2006-01-02", v, time.Local)
			if err == nil {
				u.Expire = t.Unix()
			}
		}
		delete(tmp, "expire")
	}
	if vv, has := tmp["hide_credential"]; has {
		v, ok := vv.(bool)
		if ok {
			u.HideCredential = v
		}
		delete(tmp, "hide_credential")
	}
	pattern := make(map[string]string)
	for key, value := range tmp {
		v, ok := value.(string)
		if ok {
			pattern[key] = v
		}
	}
	u.Pattern = pattern
	return nil
}

func dealAPPAppend(data map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range data {
		if key == "additional" || key == "auth" {
			vs, ok := value.([]interface{})
			if !ok {
				continue
			}
			params := make([]interface{}, 0, len(vs))
			for _, v := range vs {
				val, ok := v.(map[string]interface{})
				if !ok {
					continue
				}
				vv, ok := val["config"]
				if ok {
					if key == "auth" {
						tmp, err := json.Marshal(vv)
						if err != nil {
							continue
						}
						var cfg Config
						err = json.Unmarshal(tmp, &cfg)
						if err != nil {
							continue
						}
						params = append(params, &cfg)
						continue
					}
					params = append(params, vv)
				}
			}
			result[key] = params
			continue
		}
		result[key] = value
	}
	return result
}
