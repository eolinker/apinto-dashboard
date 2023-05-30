package service

import "github.com/eolinker/eosc/common/bean"

func init() {
	remoteStorage := newRemoteStorageService()
	bean.Injection(&remoteStorage)

}
