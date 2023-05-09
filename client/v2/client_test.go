package v2

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestClient(t *testing.T) {
	client := NewClient("http://127.0.0.1:9400")
	result, err := client.List("router")
	if err != nil {
		t.Error(err)
		return
	}
	for _, r := range result {
		fmt.Println(r)
	}
	data, err := json.Marshal(result)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(string(data))
}

func TestStructKey(t *testing.T) {
	tmp := BasicInfo{}
	typ := reflect.TypeOf(tmp)
	for i := 0; i < typ.NumField(); i++ {
		p := typ.Field(i)
		tag := p.Tag.Get("json")
		fmt.Println(tag)
	}
}
