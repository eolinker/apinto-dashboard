package common

import (
	"errors"
	"net/url"
	"os"
	"path/filepath"
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

func GetCurrentPath() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}
	currentPath, err := filepath.EvalSymlinks(filepath.Dir(execPath))
	if err != nil {
		return "", err
	}
	return currentPath, nil
}
