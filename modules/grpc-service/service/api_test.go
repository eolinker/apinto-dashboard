package service

import (
	"strings"
	"testing"
)

func TestReplacePath(t *testing.T) {
	path := "/test/abc/{param1}/{param2}"
	requestPath := "/test/*"
	index := strings.Index(requestPath, "*")
	proxyPath := "/api/"
	if index != -1 {
		if index > 0 {
			path = strings.Replace(path, requestPath[:index], proxyPath, 1)
		}
	}
	t.Log(path)
}
