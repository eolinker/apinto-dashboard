package service

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/remote_storage"
	"github.com/eolinker/apinto-dashboard/modules/remote_storage/entry"
	"github.com/eolinker/apinto-dashboard/modules/remote_storage/model"
	"github.com/eolinker/apinto-dashboard/modules/remote_storage/store"
	"github.com/eolinker/eosc/common/bean"
	"gorm.io/gorm"
)

type remoteStorageService struct {
	store store.IRemoteStorage
}

func (c *remoteStorageService) Get(module, key string) (*model.RemoteStorage, error) {
	ent, err := c.store.Get(context.Background(), module, key)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &model.RemoteStorage{
				Module: module,
				Key:    key,
				Object: "",
			}, nil
		}
		return nil, err
	}
	return &model.RemoteStorage{
		Module:   ent.Module,
		Key:      ent.Key,
		Object:   ent.Object,
		UpdateAt: ent.UpdatedAt,
		CreateAt: ent.CreatedAt,
	}, nil
}

func (c *remoteStorageService) Save(module, key string, object interface{}) error {
	info := &entry.RemoteKeyObject{
		Module: module,
		Key:    key,
		Object: object,
	}
	return c.store.Save(context.Background(), info)
}

func newRemoteStorageService() remote_storage.IRemoteStorageService {
	c := &remoteStorageService{}
	bean.Autowired(&c.store)
	return c
}
