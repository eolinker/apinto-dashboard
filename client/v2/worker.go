package v2

import (
	"encoding/json"
	"reflect"
)

var existKeys = getFieldTags(reflect.TypeOf(BasicInfo{}))

type WorkerInfo[T any] struct {
	BasicInfo *T
	Append    map[string]interface{}
}

func (r *WorkerInfo[T]) UnmarshalJSON(bytes []byte) error {
	basicInfo := new(T)
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

func (r *WorkerInfo[T]) MarshalJSON() ([]byte, error) {
	tmp := make(map[string]interface{})
	for key, value := range r.Append {
		tmp[key] = value
	}
	val := reflect.ValueOf(r.BasicInfo)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	num := val.NumField()
	typ := val.Type()
	for i := 0; i < num; i++ {
		tag := typ.Field(i).Tag.Get("json")
		if tag == "-" || tag == "" {
			continue
		}
		tmp[tag] = val.Field(i).Interface()
	}
	return json.Marshal(tmp)
}

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
