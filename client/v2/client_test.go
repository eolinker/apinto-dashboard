package v2

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestClient(t *testing.T) {
	client := NewClient("http://172.18.189.43:31194")
	result, err := client.List("service")
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

func TestVariable(t *testing.T) {
	var data = `{"config":{"address":["${consul_addr}"],"address_variable":["${consul_addr}"],"use_variable":[1]},"create":"","description":"","driver":"consul","id":"","name":"test_consul","profession":"discovery","update":"","version":"20230510230425"}`
	client := NewClient("http://127.0.0.1:9400")
	var info WorkerInfo[BasicInfo]
	json.Unmarshal([]byte(data), &info)
	err := client.Set("discovery", "aaa", &info)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestSetClient(t *testing.T) {
	var data = `{"config":{"address":["${consul_addr}"],"address_variable":["${consul_addr}"],"use_variable":[1]},"create":"","description":"","driver":"consul","id":"","name":"test_consul","profession":"discovery","update":"","version":"20230510230425"}`
	client := NewClient("http://172.18.189.43:31194")
	var info WorkerInfo[BasicInfo]
	json.Unmarshal([]byte(data), &info)
	err := client.Set("service", "aaa", &info)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestVersionClient(t *testing.T) {
	client := NewClient("http://172.18.189.43:31194")
	result, err := client.Versions("service")
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
