package store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/entry"
	"gorm.io/gorm"
)

type IServiceId interface {
	SetVersionId(int)
}

type IRuntimeId interface {
	SetRuntimeId(int)
}

type BaseKindHandler[T any, K any] interface {
	Encode(t *T) *K
	Decode(*K) *T
	Kind() string
}

type BaseKindStore[T any, K any] struct {
	*BaseStore[K]
	BaseKindHandler[T, K]
}

func (k *BaseKindStore[T, K]) Get(ctx context.Context, id int) (*T, error) {
	var ii interface{} = new(K)
	isTarget := false
	if _, ok := ii.(*entry.Stat); ok {
		isTarget = ok
	}
	if _, ok := ii.(*entry.Runtime); ok {
		isTarget = ok
	}
	if isTarget {
		vd, err := k.BaseStore.FirstQuery(ctx, "`target` = ? and `kind` = ?", []interface{}{id, k.Kind()}, "")
		if err != nil {
			return nil, err
		}
		return k.decode(vd), nil
	}
	vd, err := k.BaseStore.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return k.decode(vd), nil
}

func (k *BaseKindStore[T, K]) Save(ctx context.Context, t *T) error {

	encode := k.encode(t)
	if err := k.BaseStore.Save(ctx, encode); err != nil {
		return err
	}
	var e interface{} = encode

	id := 0
	if table, ok := e.(Table); ok {
		id = table.IdValue()
	}

	var v interface{} = t
	if table, ok := v.(IServiceId); ok {
		table.SetVersionId(id)
	}

	if table, ok := v.(IRuntimeId); ok {
		table.SetRuntimeId(id)
	}

	return nil
}

func (k *BaseKindStore[T, K]) Insert(ctx context.Context, ts ...*T) error {
	vds := k.encodes(ts)

	return k.BaseStore.Insert(ctx, vds...)
}

func (k *BaseKindStore[T, K]) List(ctx context.Context, m map[string]interface{}) ([]*T, error) {
	list, err := k.BaseStore.List(ctx, m)
	if err != nil {
		return nil, err
	}
	return k.decodes(list), err
}

func (k *BaseKindStore[T, K]) ListQuery(ctx context.Context, sql string, args []interface{}, order string) ([]*T, error) {
	list, err := k.BaseStore.ListQuery(ctx, sql, args, order)
	if err != nil {
		return nil, err
	}
	return k.decodes(list), err
}

func (k *BaseKindStore[T, K]) First(ctx context.Context, m map[string]interface{}) (*T, error) {
	query, err := k.BaseStore.First(ctx, m)
	if err != nil {
		return nil, err
	}
	return k.decode(query), nil
}

func (k *BaseKindStore[T, K]) FirstQuery(ctx context.Context, sql string, args []interface{}, order string) (*T, error) {
	query, err := k.BaseStore.FirstQuery(ctx, sql, args, order)
	if err != nil {
		return nil, err
	}
	return k.decode(query), nil
}

func (k *BaseKindStore[T, K]) ListPage(ctx context.Context, sql string, pageNum, pageSize int, args []interface{}, order string) ([]*T, int, error) {
	list, i, err := k.BaseStore.ListPage(ctx, sql, pageNum, pageSize, args, order)
	if err != nil {
		return nil, 0, err
	}
	return k.decodes(list), i, err
}

func (k *BaseKindStore[T, K]) UpdateByUnique(ctx context.Context, t *T, fields []string) error {
	return k.BaseStore.UpdateByUnique(ctx, k.encode(t), fields)
}

func (k *BaseKindStore[T, K]) UpdateWhere(ctx context.Context, t *T, m map[string]interface{}) (int, error) {
	return k.BaseStore.UpdateWhere(ctx, k.encode(t), m)

}

func (k *BaseKindStore[T, K]) Update(ctx context.Context, t *T) (int, error) {
	return k.BaseStore.Update(ctx, k.encode(t))
}

func (k *BaseKindStore[T, K]) DB(ctx context.Context) *gorm.DB {
	return k.BaseStore.DB(ctx)
}

func (k *BaseKindStore[T, K]) decode(i *K) *T {
	return k.Decode(i)
}

func (k *BaseKindStore[T, K]) decodes(list []*K) []*T {
	vds := make([]*T, 0, len(list))
	for _, l := range list {
		vds = append(vds, k.Decode(l))
	}
	return vds
}

func (k *BaseKindStore[T, K]) encode(i *T) *K {
	return k.Encode(i)
}

func (k *BaseKindStore[T, K]) encodes(list []*T) []*K {
	vds := make([]*K, 0, len(list))
	for _, l := range list {
		vds = append(vds, k.Encode(l))
	}
	return vds
}

func CreateBaseKindStore[T any, K any](t BaseKindHandler[T, K], db IDB) *BaseKindStore[T, K] {
	k := &BaseKindStore[T, K]{
		BaseKindHandler: t,
		BaseStore:       CreateStore[K](db),
	}
	return k
}
