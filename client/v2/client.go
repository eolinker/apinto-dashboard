package v2

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type IClient interface {
	List(profession string) ([]*WorkerInfo[BasicInfo], error)
	Versions(profession string, drivers ...string) (map[string]string, error)
	Version(profession string, name string) (string, error)
	Info(profession string, name string) (*WorkerInfo[BasicInfo], error)
	Set(profession string, name string, info *WorkerInfo[BasicInfo]) error
	Delete(profession string, name string) error
	Ping() error
	Cluster() (*Cluster, error)
	Addr() string
}

func NewClient(addr string) IClient {
	client := http.Client{
		Transport: http.DefaultTransport,
		Timeout:   10 * time.Second,
	}
	return &Client{
		addr:   strings.TrimSuffix(addr, "/"),
		client: &client,
	}
}

type Client struct {
	addr   string
	client *http.Client
}

func (c *Client) Version(profession string, name string) (string, error) {
	respBody, err := sendTo(c.client, http.MethodGet, fmt.Sprintf("%s/api/%s/%s", c.addr, profession, name), nil, "", []int{200})
	if err != nil {
		return "", err
	}
	result := new(BasicInfo)
	err = json.Unmarshal(respBody, result)
	if err != nil {
		return "", fmt.Errorf("unmarshal data error: %w, body is %s", err, respBody)
	}
	return result.Version, nil
}

func (c *Client) Versions(profession string, drivers ...string) (map[string]string, error) {
	respBody, err := sendTo(c.client, http.MethodGet, fmt.Sprintf("%s/api/%s", c.addr, profession), nil, "", []int{200})
	if err != nil {
		return nil, err
	}
	versionMap := make(map[string]string)
	list := make([]*BasicInfo, 0)
	err = json.Unmarshal(respBody, &list)
	if err != nil {
		return nil, err
	}
	all := len(drivers) < 1
	for _, l := range list {
		if all {
			versionMap[l.Name] = l.Version
			continue
		}
		for _, driver := range drivers {
			if l.Driver != driver {
				continue
			}
			versionMap[l.Name] = l.Version
		}
	}
	return versionMap, nil
}

func (c *Client) Ping() error {
	_, err := sendTo(c.client, http.MethodGet, fmt.Sprintf("%s/system/info", c.addr), nil, "", []int{200})
	return err
}

func (c *Client) Cluster() (*Cluster, error) {
	respBody, err := sendTo(c.client, http.MethodGet, fmt.Sprintf("%s/system/info", c.addr), nil, "", []int{200})
	if err != nil {
		return nil, err
	}
	result := new(Cluster)
	err = json.Unmarshal(respBody, result)
	if err != nil {
		return nil, fmt.Errorf("unmarshal data error: %w, body is %s", err, respBody)
	}
	return result, nil
}

func (c *Client) Addr() string {
	return c.addr
}

func (c *Client) List(profession string) ([]*WorkerInfo[BasicInfo], error) {
	respBody, err := sendTo(c.client, http.MethodGet, fmt.Sprintf("%s/api/%s", c.addr, profession), nil, "", []int{200})
	if err != nil {
		return nil, err
	}
	result := make([]*WorkerInfo[BasicInfo], 0)
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return nil, fmt.Errorf("unmarshal data error: %w, body is %s", err, respBody)
	}
	return result, nil
}

func (c *Client) Info(profession string, name string) (*WorkerInfo[BasicInfo], error) {
	respBody, err := sendTo(c.client, http.MethodGet, fmt.Sprintf("%s/api/%s/%s", c.addr, profession, name), nil, "", []int{200})
	if err != nil {
		return nil, err
	}
	result := new(WorkerInfo[BasicInfo])
	err = json.Unmarshal(respBody, result)
	if err != nil {
		return nil, fmt.Errorf("unmarshal data error: %w, body is %s", err, respBody)
	}
	return result, nil
}

func (c *Client) Set(profession string, name string, info *WorkerInfo[BasicInfo]) error {
	header := http.Header{}
	header.Set("content-type", "application/json")
	body, _ := json.Marshal(info)
	_, err := sendTo(c.client, http.MethodPost, fmt.Sprintf("%s/api/%s/%s", c.addr, profession, name), header, string(body), []int{200})
	return err
}

func (c *Client) Delete(profession string, name string) error {
	_, err := sendTo(c.client, http.MethodDelete, fmt.Sprintf("%s/api/%s/%s", c.addr, profession, name), nil, "", []int{200, 404})
	return err
}

func sendTo(client *http.Client, method string, uri string, headers http.Header, body string, successStatus []int) ([]byte, error) {
	req, err := http.NewRequest(method, uri, strings.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("get clister info error: %w", err)
	}
	req.Header = headers
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request error: %w", err)
	}
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body error: %w", err)
	}
	success := false
	for _, status := range successStatus {
		if status == resp.StatusCode {
			success = true
			break
		}
	}
	if !success {
		return nil, fmt.Errorf("error status code: %d,body is %s", resp.StatusCode, responseBody)
	}
	return responseBody, nil
}
