package dynamic_entry

import (
	"encoding/json"
	"reflect"
	"time"
)

type DynamicPublishVersion struct {
	Id          int                   `json:"id"`
	ClusterId   int                   `json:"cluster_id"`
	NamespaceId int                   `json:"namespace_id"`
	Publish     *DynamicPublishConfig `json:"publish"`
	Operator    int                   `json:"operator"`
	CreateTime  time.Time             `json:"create_time"`
}

func (v *DynamicPublishVersion) SetVersionId(id int) {
	v.Id = id
}

type DynamicPublishConfig struct {
	*BasicInfo
	Append map[string]interface{}
}

var existKeys = getFieldTags(reflect.TypeOf(BasicInfo{}))

type BasicInfo struct {
	Profession  string `json:"profession"`
	Id          string `json:"id"`
	Name        string `json:"name"`
	Driver      string `json:"driver"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Create      string `json:"create"`
	Update      string `json:"update"`
}

func getFieldTags(typ reflect.Type) []string {
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	num := typ.NumField()
	tags := make([]string, 0, num)
	for i := 0; i < num; i++ {
		tags = append(tags, typ.Field(i).Tag.Get("json"))
	}
	return tags
}

func (r *DynamicPublishConfig) UnmarshalJSON(bytes []byte) error {
	basicInfo := new(BasicInfo)
	err := json.Unmarshal(bytes, basicInfo)
	if err != nil {
		return err
	}
	tmp := make(map[string]interface{})
	err = json.Unmarshal(bytes, &tmp)
	if err != nil {
		return err
	}
	r.BasicInfo = basicInfo
	for _, key := range existKeys {
		delete(tmp, key)
	}
	r.Append = tmp
	return nil
}
