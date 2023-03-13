package store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/entry"
)

type IApplicationAuthPublishStore interface {
	IBaseStore[entry.ApplicationAuthPublish]
	Reset(ctx context.Context, clusterId, applicationId int, list []*entry.ApplicationAuthPublish) error
	DeleteByClusterIdAppId(ctx context.Context, clusterId, applicationId int) error
	GetList(ctx context.Context, clusterId, applicationId int) ([]*entry.ApplicationAuthPublish, error)
}

type applicationAuthPublishStore struct {
	*baseStore[entry.ApplicationAuthPublish]
}

func newApplicationAuthPublishStore(db IDB) IApplicationAuthPublishStore {
	return &applicationAuthPublishStore{baseStore: createStore[entry.ApplicationAuthPublish](db)}
}

// Reset 外部保证事务
func (a *applicationAuthPublishStore) Reset(ctx context.Context, clusterId, applicationId int, list []*entry.ApplicationAuthPublish) error {
	_, err := a.DeleteWhere(ctx, map[string]interface{}{"`cluster`": clusterId, "`application`": applicationId})
	if err != nil {
		return err
	}
	return a.Insert(ctx, list...)
}

// GetList
func (a *applicationAuthPublishStore) GetList(ctx context.Context, clusterId, applicationId int) ([]*entry.ApplicationAuthPublish, error) {
	return a.ListQuery(ctx, "`cluster` = ? and `application` = ?", []interface{}{clusterId, applicationId}, "")
}

func (a *applicationAuthPublishStore) DeleteByClusterIdAppId(ctx context.Context, clusterId, applicationId int) error {
	_, err := a.DeleteWhere(ctx, map[string]interface{}{"`cluster`": clusterId, "`application`": applicationId})
	return err
}
