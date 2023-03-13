package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

type APPListData struct {
	Applications []*AppItem `json:"applications"`
}

type AppItem struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type APPListResp struct {
	Code int          `json:"code"`
	Data *APPListData `json:"data"`
	Msg  string       `json:"msg"`
}

// getAppItems 从控制台app列表
func getAppItems() ([]*AppItem, error) {
	//client := http.Client{
	//	Timeout: 10 * time.Minute,
	//}
	//
	//url := fmt.Sprintf("http://%s/api/application/enum", runConfig.ApserverAddr)
	//req, err := http.NewRequest(http.MethodGet, url, strings.NewReader(""))
	//if err != nil {
	//	return nil, err
	//}
	//
	////组装请求cookie， cookie可以登陆开发环境admin用F12去获取
	//namespaceCookie := &http.Cookie{
	//	Name:  "namespace",
	//	Value: "default",
	//}
	//sessionCookie := &http.Cookie{
	//	Name:  "Session",
	//	Value: runConfig.ApserverSession,
	//}
	//req.AddCookie(namespaceCookie)
	//req.AddCookie(sessionCookie)
	//
	///*组装请求query
	//分组uuid可以从开发环境获取
	//*/
	//reqQuery := req.URL.Query()
	//reqQuery.Add("namespace", "default")
	//
	//req.URL.RawQuery = reqQuery.Encode()
	//
	//resp, err := client.Do(req)
	//if err != nil {
	//	return nil, err
	//}
	//defer resp.Body.Close()
	//body, err := io.ReadAll(resp.Body)
	//if err != nil {
	//	return nil, err
	//}
	//
	//appListResp := &APPListResp{}
	//err = json.Unmarshal(body, appListResp)
	//if err != nil {
	//	return nil, err
	//}
	//if appListResp.Code != 0 {
	//	return nil, fmt.Errorf("get service items fail. err: %s", appListResp.Msg)
	//}

	file := path.Join("./export/app", "apps.json")

	data, err := os.ReadFile(file)
	if err != nil {
		fmt.Println("File reading error", err)
		return nil, err
	}
	appListData := &APPListResp{}
	err = json.Unmarshal(data, appListData)
	if err != nil {
		return nil, err
	}

	return appListData.Data.Applications, nil
}
