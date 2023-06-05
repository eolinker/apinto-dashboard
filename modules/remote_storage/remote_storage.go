package remote_storage

import "github.com/eolinker/apinto-dashboard/modules/remote_storage/model"

type IRemoteStorageService interface {
	Get(module, key string) (*model.RemoteStorage, error)
	Save(module, key string, object interface{}) error
}
