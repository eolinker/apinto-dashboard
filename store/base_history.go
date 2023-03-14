package store

import (
	"context"
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/modules/base/history-entry"
	"time"
)

type BaseHistoryStore[T any] interface {
	History(ctx context.Context, namespace, target int, OldValue, newValue interface{}, optType history_entry.OptType, Operator int) error
	HistoryAdd(ctx context.Context, namespace, target int, newValue interface{}, Operator int) error
	HistoryEdit(ctx context.Context, namespace, target int, OldValue, newValue interface{}, Operator int) error
	HistoryDelete(ctx context.Context, namespace, target int, oldValue interface{}, Operator int) error
	List(ctx context.Context, namespace int, target ...int) ([]*T, error)
	Page(ctx context.Context, namespace int, pageNum, pageSize int, target ...int) ([]*T, int, error)
	LatestOperators(ctx context.Context, target ...int) (map[int]int, error)
	Latest(ctx context.Context, target ...int) (map[int]*T, error)
	LatestOne(ctx context.Context, target int) (*T, error)
}
type DecodeHistory[T any] interface {
	Decode(e *history_entry.History) *T
}
type BaseHistory[T any] struct {
	kind    history_entry.HistoryKind
	myStore *BaseStore[history_entry.History]
	decoder DecodeHistory[T]
}

func (b *BaseHistory[T]) History(ctx context.Context, namespace int, target int, OldValue, newValue interface{}, optType history_entry.OptType, operator int) error {

	oldValueStr, err := encodeValue(OldValue)
	if err != nil {
		return err
	}

	newValueStr, err := encodeValue(newValue)
	if err != nil {
		return err
	}
	history := &history_entry.History{
		Id:          0,
		NamespaceID: namespace,
		Kind:        b.kind,
		TargetID:    target,
		OldValue:    oldValueStr,
		NewValue:    newValueStr,
		OptType:     optType,
		Operator:    operator,
		OptTime:     time.Now(),
	}
	return b.myStore.Insert(ctx, history)
}
func (b *BaseHistory[T]) HistoryAdd(ctx context.Context, namespace, target int, newValue interface{}, Operator int) error {
	return b.History(ctx, namespace, target, nil, newValue, history_entry.OptAdd, Operator)
}
func (b *BaseHistory[T]) HistoryEdit(ctx context.Context, namespace, target int, OldValue, newValue interface{}, Operator int) error {
	return b.History(ctx, namespace, target, OldValue, newValue, history_entry.OptEdit, Operator)
}

func (b *BaseHistory[T]) HistoryDelete(ctx context.Context, namespace, target int, oldValue interface{}, Operator int) error {
	return b.History(ctx, namespace, target, oldValue, nil, history_entry.OptDelete, Operator)
}
func (b *BaseHistory[T]) List(ctx context.Context, namespace int, target ...int) ([]*T, error) {
	query := map[string]interface{}{
		"kind":      b.kind,
		"namespace": namespace,
	}
	if len(target) > 0 {
		query["target"] = target[0]
	}
	hs, err := b.myStore.List(ctx, query)
	if err != nil {
		return nil, err
	}
	list := make([]*T, 0, len(hs))
	for _, e := range hs {
		list = append(list, b.decoder.Decode(e))
	}
	return list, nil
}

func (b *BaseHistory[T]) Page(ctx context.Context, namespace int, pageNum, pageSize int, target ...int) ([]*T, int, error) {
	where := "`namespace`=? and `kind`=?"
	arg := []interface{}{namespace, b.kind}
	if len(target) > 0 {
		where += " and `target` in (?)"
		arg = append(arg, target)
	}

	page, i, err := b.myStore.ListPage(ctx, where, pageNum, pageSize, arg, "id desc")
	if err != nil {
		return nil, 0, err
	}
	rs := make([]*T, 0, len(page))
	for _, p := range page {
		rs = append(rs, b.decoder.Decode(p))
	}
	return rs, i, nil
}

type KV struct {
	Operator int `gorm:"column:operator"`
	TargetID int `gorm:"column:target"`
}

func (b *BaseHistory[T]) LatestOperators(ctx context.Context, target ...int) (map[int]int, error) {

	rs := make([]*KV, 0, len(target))
	err := b.myStore.DB(ctx).Raw("SELECT `operator`,`target` FROM `history` WHERE `id` IN( SELECT max(`id`) FROM history WHERE `kind` = ? AND `target` IN (?)   GROUP BY `target`)", b.kind, target).Scan(&rs).Error
	if err != nil {
		return nil, err
	}
	rm := make(map[int]int, len(rs))
	for _, kv := range rs {
		rm[kv.TargetID] = rm[kv.Operator]
	}
	return rm, nil
}

func (b *BaseHistory[T]) Latest(ctx context.Context, target ...int) (map[int]*T, error) {
	rs := make([]*history_entry.History, 0, len(target))

	err := b.myStore.DB(ctx).Raw("SELECT * FROM `history` WHERE `id` IN( SELECT max(`id`) FROM history WHERE `kind` = ? AND `target` IN (?)   GROUP BY `target`)").Scan(&rs).Error
	if err != nil {
		return nil, err
	}

	rm := common.SliceToMap(rs, func(t *history_entry.History) int {
		return t.TargetID
	})

	tm := make(map[int]*T, len(target))
	for t, v := range rm {
		tm[t] = b.decoder.Decode(v)
	}
	return tm, nil
}

func (b *BaseHistory[T]) LatestOne(ctx context.Context, target int) (*T, error) {
	first, err := b.myStore.FirstQuery(ctx, "`kind` = ? and `target` = ?", []interface{}{b.kind, target}, "id desc")
	if err != nil {
		return nil, err
	}
	return b.decoder.Decode(first), nil
}

func CreateHistory[T any](decoder DecodeHistory[T], db IDB, kind history_entry.HistoryKind) BaseHistoryStore[T] {
	return &BaseHistory[T]{
		kind:    kind,
		myStore: CreateStore[history_entry.History](db),
		decoder: decoder,
	}
}

func encodeValue(value interface{}) (string, error) {
	if value == nil {
		return "", nil
	}
	switch v := value.(type) {
	case string:
		return v, nil
	case *string:
		return *v, nil

	}
	bs, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}
