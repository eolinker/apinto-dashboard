package v1

import (
	"bytes"
	"fmt"
	"github.com/eolinker/eosc/log"
	"io"
	"net/http"
	"time"
)

func requestDo(method string, url string, body []byte) ([]byte, int, error) {
	req, err := newRequest(method, url, body)
	if err != nil {
		log.Error("new request:", err)
		return nil, 500, err
	}
	client := http.Client{Timeout: time.Second * 3}
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Error("do request:", err)
		return nil, 500, err
	}
	data, err := readBody(resp.Body)
	if err != nil {
		log.Error("read body:", err)

		return nil, 500, err
	}
	return data, resp.StatusCode, nil
}

func newRequest(method string, url string, body []byte) (*http.Request, error) {
	if body == nil {
		return http.NewRequest(method, url, nil)
	}
	return http.NewRequest(method, url, bytes.NewReader(body))
}

func readBody(r io.ReadCloser) ([]byte, error) {
	defer func() {
		if err := r.Close(); err != nil {
			log.Errorf("failed to close body, err: %s", err.Error())
		}
	}()
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read body, err: %s", err.Error())
	}
	return data, nil
}
