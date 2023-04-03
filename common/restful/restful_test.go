package restful

import (
	"encoding/json"
	"fmt"
	_ "github.com/stretchr/testify/assert"
)

var (
	config = BuildConfig(nil, nil, "http://demo.apinto.com:8280/")
)

func ExampleRPC() {
	type Demo struct {
		Data map[string]string `json:"data,omitempty"`
	}

	demo := Rpc[Demo, Demo](config, "post", "/demo")
	input := &Demo{Data: map[string]string{"a": "a", "B": "b"}}
	response, err := demo.Build().Request(input)
	if err != nil {
		return
	}
	marshal, err := json.Marshal(response.Data)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(marshal))
	//output: {}
}
func ExampleSimple() {

	demo := Simple(config, "post", "/demo/{x}")
	//input := Demo{Data: map[string]string{"a": "a", "B": "b"}}
	response, err := demo.Build().Request("x_value")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response)
	//output: ""
}
func ExampleOneWay() {
	type Demo struct {
		Data map[string]string `json:"data,omitempty"`
	}

	demo := OneWay[Demo](config, "post", "/demo/{x}")
	input := &Demo{Data: map[string]string{"a": "a", "B": "b"}}
	response, err := demo.Build().Request(input, "x_value")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response)
	//output: ""
}
func ExampleCall() {
	type Demo map[string]string

	demo := Call[Demo](config, "post", "/demo/{x}")

	response, err := demo.Build().Request("x_value")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response)
	//output: ""
}
