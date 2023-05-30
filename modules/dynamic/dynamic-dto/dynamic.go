package dynamic_dto

import (
	"encoding/json"
	"reflect"
)

var existKeys = getFieldTags(reflect.TypeOf(BasicInfo{}))

type WorkerInfo struct {
	*BasicInfo
	Append map[string]interface{}
}

type BasicInfo struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Driver      string `json:"driver"`
	Description string `json:"description"`
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

func (r *WorkerInfo) UnmarshalJSON(bytes []byte) error {
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
