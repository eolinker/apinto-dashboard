package common

import (
	"errors"
	"strings"
)

var (
	ClusterNotExist = errors.New("cluster does not exist")
)

const workerNotExistErrStr = "worker-data not exits"

// CheckWorkerNotExist 检查报错信息是否为worker实例不存在
func CheckWorkerNotExist(err error) error {
	if err != nil && strings.Contains(err.Error(), workerNotExistErrStr) {
		return nil
	}
	return err
}
