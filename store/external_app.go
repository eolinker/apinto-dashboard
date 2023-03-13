package store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/entry"
	"time"
)

var _ IExternalApplicationStore = (*externalApplicationStore)(nil)

type IExternalApplicationStore interface {
	IBaseStore[entry.ExternalApplication]
	GetByUUID(ctx context.Context, namespaceId int, uuid string) (*entry.ExternalApplication, error)
	GetList(ctx context.Context, namespaceId int) ([]*entry.ExternalApplication, error)
	SoftDelete(ctx context.Context, operator, id int) error
	FlushToken(ctx context.Context, operator, id int, token string) error
	GetByToken(ctx context.Context, namespaceId int, token string) (*entry.ExternalApplication, error)
}

type externalApplicationStore struct {
	*BaseStore[entry.ExternalApplication]
}

func (e *externalApplicationStore) GetByUUID(ctx context.Context, namespaceId int, uuid string) (*entry.ExternalApplication, error) {
	return e.FirstQuery(ctx, "`namespace` = ? and `uuid` = ?", []interface{}{namespaceId, uuid}, "")
}

func (e *externalApplicationStore) GetList(ctx context.Context, namespaceId int) ([]*entry.ExternalApplication, error) {
	return e.ListQuery(ctx, "`namespace` = ? and `is_delete` = 0 ", []interface{}{namespaceId}, "update_time desc")
}

func (e *externalApplicationStore) SoftDelete(ctx context.Context, operator, id int) error {
	_, err := e.UpdateWhere(ctx, &entry.ExternalApplication{Id: id}, map[string]interface{}{"is_delete": true, "operator": operator})
	return err
}

func (e *externalApplicationStore) FlushToken(ctx context.Context, operator, id int, token string) error {
	_, err := e.UpdateWhere(ctx, &entry.ExternalApplication{Id: id}, map[string]interface{}{"token": token, "operator": operator, "update_time": time.Now()})
	return err
}
func (e *externalApplicationStore) GetByToken(ctx context.Context, namespaceId int, token string) (*entry.ExternalApplication, error) {
	return e.FirstQuery(ctx, "`namespace` = ? and `token` = ?", []interface{}{namespaceId, token}, "")
}

func newExternalApplicationStore(db IDB) IExternalApplicationStore {
	return &externalApplicationStore{BaseStore: CreateStore[entry.ExternalApplication](db)}
}
