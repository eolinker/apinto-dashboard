package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

const (
	fileDir  = "./export/api"
	fileName = "items.json"
)

type APIListItem struct {
	APIUUID     string   `json:"uuid"`
	APIName     string   `json:"name"`
	Method      []string `json:"method"`
	ServiceName string   `json:"service"`
	RequestPath string   `json:"request_path"`
}

type APIListData struct {
	Apis  []*APIListItem `json:"apis"`
	Total int            `json:"total"`
}

type APIListResp struct {
	Code int          `json:"code"`
	Data *APIListData `json:"data"`
	Msg  string       `json:"msg"`
}

// getApiInfoAndSave 从控制台获取api信息,并存入文件
func getApiInfoAndSave() error {
	client := http.Client{
		Timeout: 10 * time.Minute,
	}

	url := fmt.Sprintf("http://%s/api/routers", runConfig.ApserverAddr)
	req, err := http.NewRequest(http.MethodGet, url, strings.NewReader(""))
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

	/*组装请求query
	分组uuid可以从开发环境获取
	*/
	reqQuery := req.URL.Query()
	reqQuery.Add("namespace", "default")
	reqQuery.Add("group_uuid", runConfig.groupUUID)
	reqQuery.Add("page_num", "1")
	reqQuery.Add("page_size", "10000")
	req.URL.RawQuery = reqQuery.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	apiListResp := &APIListResp{}
	err = json.Unmarshal(body, apiListResp)
	if err != nil {
		return err
	}
	if apiListResp.Code != 0 {
		return fmt.Errorf("get api items fail. err: %s", apiListResp.Msg)
	}

	return genFile(apiListResp.Data)
}

func genFile(data *APIListData) error {
	fileData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = os.MkdirAll(fileDir, 0777)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path.Join(fileDir, fileName), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0777)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	n, err := writer.Write(fileData)
	writer.Flush()
	if err != nil {
		fmt.Println("写入文件失败。 err: ", err)
		return err
	} else {
		fmt.Println("写入文件成功。 写入常茹: ", n)
	}

	return err
}

func getApiInfosFromFile() ([]*APIListItem, error) {
	file := path.Join(fileDir, fileName)

	data, err := os.ReadFile(file)
	if err != nil {
		fmt.Println("File reading error", err)
		return nil, err
	}
	apiListData := &APIListData{}
	err = json.Unmarshal(data, apiListData)
	if err != nil {
		return nil, err
	}

	return apiListData.Apis, nil
}
