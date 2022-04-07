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

func (a *adminClient) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Content-Type", "application/json")
	return a.client.Do(req)
}

func (a *adminClient) setUrl(url string) {
	a.baseUrl = url
}

func (a *adminClient) List(profession string) (*response, error) {
	url := fmt.Sprintf("%s/%s", a.baseUrl, profession)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := a.do(req)
	if err != nil {
		return nil, err
	}
	data, err := readBody(resp.Body)
	if err != nil {
		return nil, err
	}
	return &response{
		data: data,
		code: resp.StatusCode,
	}, nil
}

func (a *adminClient) Get(profession string, name string) (*response, error) {
	url := fmt.Sprintf("%s/%s/%s", a.baseUrl, profession, name)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := a.do(req)
	if err != nil {
		return nil, err
	}
	data, err := readBody(resp.Body)
	if err != nil {
		return nil, err
	}
	return &response{
		data: data,
		code: resp.StatusCode,
	}, nil
}

func (a *adminClient) Update(profession string, name string, body []byte) (*response, error) {
	url := fmt.Sprintf("%s/%s/%s", a.baseUrl, profession, name)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	resp, err := a.do(req)
	if err != nil {
		return nil, err
	}
	data, err := readBody(resp.Body)
	if err != nil {
		return nil, err
	}
	return &response{
		data: data,
		code: resp.StatusCode,
	}, nil
}

func (a *adminClient) Create(profession string, body []byte) (*response, error) {
	url := fmt.Sprintf("%s/%s", a.baseUrl, profession)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := a.do(req)
	if err != nil {
		return nil, err
	}
	data, err := readBody(resp.Body)
	if err != nil {
		return nil, err
	}
	return &response{
		data: data,
		code: resp.StatusCode,
	}, nil
}

func (a *adminClient) Delete(profession string, name string) (*response, error) {
	url := fmt.Sprintf("%s/%s/%s", a.baseUrl, profession, name)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := a.do(req)
	if err != nil {
		return nil, err
	}
	data, err := readBody(resp.Body)
	if err != nil {
		return nil, err
	}
	return &response{
		data: data,
		code: resp.StatusCode,
	}, nil
}

func (a *adminClient) Enable(profession string, name string) (*response, error) {
	// Todo 等apinto完成对接
	return &response{
		data: []byte("implementing"),
		code: http.StatusNotImplemented,
	}, nil
}

func (a *adminClient) Patch(profession string, name string, body []byte) (*response, error) {
	// Todo 等apinto完成对接
	url := fmt.Sprintf("%s/%s/%s", a.baseUrl, profession, name)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	resp, err := a.do(req)
	if err != nil {
		return nil, err
	}
	data, err := readBody(resp.Body)
	if err != nil {
		return nil, err
	}
	return &response{
		data: data,
		code: resp.StatusCode,
	}, nil
}

func (a *adminClient) Render(profession string, driver string) (interface{}, error) {
	url := fmt.Sprintf("%s/%s/_render/%s", a.baseUrl, profession, driver)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := a.do(req)
	if err != nil {
		return nil, err
	}
	data, err := readBody(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
