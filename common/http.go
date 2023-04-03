package common

import (
	"io"
	"net/http"
	"strings"
)

// NewRequest body = strings.NewReader(string(by))
func NewRequest(method, url string, body io.Reader) (by []byte, err error) {
	method = strings.ToUpper(method)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, body)
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	defer req.Body.Close()
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
