package apinto

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var (
	nodes  = "http://127.0.0.1:9400/api"
	client = NewAdminClient()
)

func Client() IClient {
	return client
}

func ResetNode(url string) error {
	if url == "" {
		return errors.New("empty node url")
	}
	url = strings.TrimSuffix(url, "/")
	client.setUrl(url)
	nodes = url
	return nil
}

type adminClient struct {
	baseUrl string
	client  *http.Client
}

func NewAdminClient() *adminClient {
	return &adminClient{
		baseUrl: nodes,
		client:  http.DefaultClient,
	}
}

func (a *adminClient) do(req *http.Request) ([]byte, int, error) {
	req.Header.Set("Content-Type", "application/json")
	resp, err := a.client.Do(req)
	if err != nil {
		return nil, 500, err
	}
	data, err := ReadBody(resp.Body)
	if err != nil {
		return nil, 500, err
	}
	return data, resp.StatusCode, nil

}

func (a *adminClient) setUrl(url string) {
	a.baseUrl = url
}

func (a *adminClient) List(profession string) ([]byte, int, error) {
	url := fmt.Sprintf("%s/%s", a.baseUrl, profession)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, err
	}
	return a.do(req)
}

func (a *adminClient) Get(profession string, name string) ([]byte, int, error) {
	url := fmt.Sprintf("%s/%s/%s", a.baseUrl, profession, name)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, err
	}
	return a.do(req)
}

func (a *adminClient) Update(profession string, name string, body []byte) ([]byte, int, error) {
	url := fmt.Sprintf("%s/%s/%s", a.baseUrl, profession, name)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return nil, 0, err
	}
	return a.do(req)
}

func (a *adminClient) Create(profession string, body []byte) ([]byte, int, error) {
	url := fmt.Sprintf("%s/%s", a.baseUrl, profession)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, 0, err
	}

	return a.do(req)
}

func (a *adminClient) Delete(profession string, name string) ([]byte, int, error) {
	url := fmt.Sprintf("%s/%s/%s", a.baseUrl, profession, name)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, 0, err
	}
	return a.do(req)
}

func (a *adminClient) Patch(profession string, name string, body []byte) ([]byte, int, error) {
	// Todo 等apinto完成对接
	url := fmt.Sprintf("%s/%s/%s", a.baseUrl, profession, name)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, 0, err
	}
	return a.do(req)
}
func (a *adminClient) PatchPath(profession string, name string, path string, body []byte) ([]byte, int, error) {
	// Todo 等apinto完成对接
	url := fmt.Sprintf("%s/%s/%s/%s", a.baseUrl, profession, name, path)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, 0, err
	}
	return a.do(req)
}

func (a *adminClient) Render(profession string, driver string) ([]byte, int, error) {
	url := fmt.Sprintf("%s/%s/_render/%s", a.baseUrl, profession, driver)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, err
	}
	return a.do(req)
}
