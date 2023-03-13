package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

const (
	createAPIJson = `{
    "group_uuid":"%s",
    "name":"%s",
    "desc":"",
    "request_path":"/test-influxdb/%s",
    "service":"%s",
    "proxy_path":"/test-influxdb/%s",
    "timeout":10000,
    "retry":0,
    "uuid":"",
    "method":[
        "POST"
    ],
    "proxy_header":[

    ],
    "match":[

    ]
}`
	apiNamePrefix = "influxdb_api_%d"
)

type ServiceListResp struct {
	Code int              `json:"code"`
	Data *ServiceItemData `json:"data"`
	Msg  string           `json:"msg"`
}

type ServiceItemData struct {
	List []string `json:"list"`
}

func createAPIBatch(apiNum int) error {
	client := &http.Client{
		Timeout: 10 * time.Minute,
	}
	//获取服务列表
	serviceItems, err := getServiceItems(client)
	if err != nil {
		return err
	}
	serviceCount := len(serviceItems)
	if serviceCount == 0 {
		return fmt.Errorf("service items can't be null. ")
	}
	url := fmt.Sprintf("http://%s/api/router", runConfig.ApserverAddr)
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

	//开发环境下专门存放influxDB测试数据api的分组uuid
	groupUUID := "ca54f3b7-54ca-43a7-b02b-278a377fa1c6"

	/*组装请求query， 需要先去删掉本地/api/routers接口对分页参数的校验，这样就可以获取某个目录下所有api
	分组uuid可以从开发环境获取
	*/
	reqQuery := req.URL.Query()
	reqQuery.Add("namespace", "default")
	req.URL.RawQuery = reqQuery.Encode()

	rand.Seed(time.Now().UnixNano())
	for i := 1; i <= apiNum; i++ {
		apiName := fmt.Sprintf(apiNamePrefix, i)
		serviceName := serviceItems[rand.Intn(serviceCount)]

		reader := strings.NewReader(fmt.Sprintf(createAPIJson, groupUUID, apiName, apiName, serviceName, apiName))
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
			fmt.Printf("create api %s success. \n", apiName)
		} else {
			fmt.Printf("create api %s fail. err: %s\n", apiName, createResp.Msg)
		}
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}

// getServiceItems 从控制台service列表
func getServiceItems(client *http.Client) ([]string, error) {
	url := fmt.Sprintf("http://%s/api/service/enum", runConfig.ApserverAddr)
	req, err := http.NewRequest(http.MethodGet, url, strings.NewReader(""))
	if err != nil {
		return nil, err
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

	/*组装请求query
	分组uuid可以从开发环境获取
	*/
	reqQuery := req.URL.Query()
	reqQuery.Add("namespace", "default")

	req.URL.RawQuery = reqQuery.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	serviceListResp := &ServiceListResp{}
	err = json.Unmarshal(body, serviceListResp)
	if err != nil {
		return nil, err
	}
	if serviceListResp.Code != 0 {
		return nil, fmt.Errorf("get service items fail. err: %s", serviceListResp.Msg)
	}

	return serviceListResp.Data.List, nil
}
