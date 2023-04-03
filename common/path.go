package common

import (
	"errors"
	"net/url"
	"strings"
)

func CheckAndFormatPath(path string) (string, error) {
	decodedPath, err := url.PathUnescape(path)
	if err != nil {
		return path, errors.New("Unescape fail. ")
	}
	if CheckPathContainsIPPort(decodedPath) {
		return path, errors.New("path is illegal. ")
	}
	if idx := strings.Index(decodedPath, "?"); idx != -1 {
		return path, errors.New("can't contain ?. ")
	}
	return decodedPath, nil
}
