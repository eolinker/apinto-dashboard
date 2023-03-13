package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	serviceJson = `{
    "name":"%s",
    "desc":"",
    "scheme":"HTTP",
    "balance":"round-robin",
    "discovery_name":"static",
    "timeout":100,
    "config":{
        "addrs_variable":"",
        "use_variable":false,
        "service_name":"",
        "static_conf":[
            {
                "weight":10,
                "addr":"demo.apinto.com:8280"
            }
        ]
    }
}`
	serviceNamePrefix = "influxdb_service_%d"
)

type CreateResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func createService(serviceNum int) error {
	client := http.Client{
		Timeout: 10 * time.Minute,
	}

	url := fmt.Sprintf("http://%s/api/service", runConfig.ApserverAddr)
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(""))
	if err != nil {
		return err
	}

	//组装请求cookie， cookie可以登陆开发环境admin用F12去获取
	namespaceCookie := &http.Cookie{
		Name:  "namespace",
		Value: "default",
	}
	sessionCookie := &http.Cookie{
		Name:  "Session",
		Value: runConfig.ApserverSession,
	}
	req.AddCookie(namespaceCookie)
	req.AddCookie(sessionCookie)

	req.Header.Add("Content-Type", "application/json")

	reqQuery := req.URL.Query()
	reqQuery.Add("namespace", "default")
	req.URL.RawQuery = reqQuery.Encode()

	for i := 1; i <= serviceNum; i++ {
		serviceName := fmt.Sprintf(serviceNamePrefix, i)
		reader := strings.NewReader(fmt.Sprintf(serviceJson, serviceName))
		req.Body = io.NopCloser(reader)

		resp, err := client.Do(req)
		if err != nil {
			return err
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		createResp := &CreateResp{}
		err = json.Unmarshal(body, createResp)
		if err != nil {
			return err
		}
		if createResp.Code == 0 {
			fmt.Printf("create service %s success. \n", serviceName)
		} else {
			fmt.Printf("create service %s fail. err: %s\n", serviceName, createResp.Msg)
		}
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}
