package application_store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/entry/application-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IApplicationAuthPublishStore interface {
	store.IBaseStore[application_entry.ApplicationAuthPublish]
	Reset(ctx context.Context, clusterId, applicationId int, list []*application_entry.ApplicationAuthPublish) error
	DeleteByClusterIdAppId(ctx context.Context, clusterId, applicationId int) error
	GetList(ctx context.Context, clusterId, applicationId int) ([]*application_entry.ApplicationAuthPublish, error)
}

type applicationAuthPublishStore struct {
	*store.BaseStore[application_entry.ApplicationAuthPublish]
}

func newApplicationAuthPublishStore(db store.IDB) IApplicationAuthPublishStore {
	return &applicationAuthPublishStore{BaseStore: store.CreateStore[application_entry.ApplicationAuthPublish](db)}
}

// Reset 外部保证事务
func (a *applicationAuthPublishStore) Reset(ctx context.Context, clusterId, applicationId int, list []*application_entry.ApplicationAuthPublish) error {
	_, err := a.DeleteWhere(ctx, map[string]interface{}{"`cluster`": clusterId, "`application`": applicationId})
	if err != nil {
		return err
	}
	return a.Insert(ctx, list...)
}

// GetList
func (a *applicationAuthPublishStore) GetList(ctx context.Context, clusterId, applicationId int) ([]*application_entry.ApplicationAuthPublish, error) {
	return a.ListQuery(ctx, "`cluster` = ? and `application` = ?", []interface{}{clusterId, applicationId}, "")
}

func (a *applicationAuthPublishStore) DeleteByClusterIdAppId(ctx context.Context, clusterId, applicationId int) error {
	_, err := a.DeleteWhere(ctx, map[string]interface{}{"`cluster`": clusterId, "`application`": applicationId})
	return err
}
