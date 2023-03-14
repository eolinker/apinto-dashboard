package quote_store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/entry/quote-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

var _ IQuoteStore = (*quoteStore)(nil)

type IQuoteStore interface {
	store.IBaseStore[quote_entry.Quote]
	DelBySource(ctx context.Context, source int, kind quote_entry.QuoteKindType) error                                                       //删除source所有的引用，比如source=discoveryID,则删除该discoveryID引用的所有东西
	DelByTarget(ctx context.Context, target int, kind quote_entry.QuoteTargetKindType) error                                                 //删除target所有的引用
	Set(ctx context.Context, source int, kind quote_entry.QuoteKindType, targetMaps map[quote_entry.QuoteTargetKindType][]int) error         //重置source的引用列表， targetMaps map[targetKind][]target
	Count(ctx context.Context, target int, targetKind quote_entry.QuoteTargetKindType) (int, error)                                          //获取target被引用的次数
	GetSourceQuote(ctx context.Context, source int, kind quote_entry.QuoteKindType) (map[quote_entry.QuoteTargetKindType][]int, error)       //获取source引用了哪些， map[targetKind][]target
	GetTargetQuote(ctx context.Context, target int, targetKind quote_entry.QuoteTargetKindType) (map[quote_entry.QuoteKindType][]int, error) //获取target 被哪些引用了， map[QuoteKindType][]source
}

type quoteStore struct {
	*store.BaseStore[quote_entry.Quote]
}

func (q *quoteStore) DelBySource(ctx context.Context, source int, kind quote_entry.QuoteKindType) error {
	_, err := q.DeleteWhere(ctx, map[string]interface{}{"`source`": source, "`kind`": kind})
	return err
}

func (q *quoteStore) DelByTarget(ctx context.Context, target int, kind quote_entry.QuoteTargetKindType) error {
	_, err := q.DeleteWhere(ctx, map[string]interface{}{"`target`": target, "`target_kind`": kind})
	return err
}

// Set 调用方开启事务
func (q *quoteStore) Set(ctx context.Context, source int, kind quote_entry.QuoteKindType, targetMaps map[quote_entry.QuoteTargetKindType][]int) error {
	if err := q.DB(ctx).Delete(q.targetType, "`source` = ? and `kind` = ?", source, kind).Error; err != nil {
		return err
	}
	list := make([]*quote_entry.Quote, 0)
	for targetKind, targets := range targetMaps {
		for _, target := range targets {
			list = append(list, &quote_entry.Quote{
				Kind:       kind,
				Source:     source,
				Target:     target,
				TargetKind: targetKind,
			})
		}
	}
	if len(list) > 0 {
		return q.DB(ctx).Create(list).Error
	}
	return nil
}

func (q *quoteStore) Count(ctx context.Context, target int, targetKind quote_entry.QuoteTargetKindType) (int, error) {
	db := q.DB(ctx)
	count := int64(0)
	if err := db.Model(q.targetType).Where("`target` = ? and `target_kind` = ?", target, targetKind).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (q *quoteStore) GetSourceQuote(ctx context.Context, source int, kind quote_entry.QuoteKindType) (map[quote_entry.QuoteTargetKindType][]int, error) {
	list, err := q.ListQuery(ctx, "`source` = ? and `kind` = ?", []interface{}{source, kind}, "")
	if err != nil {
		return nil, err
	}
	resMap := make(map[quote_entry.QuoteTargetKindType][]int)
	for _, quote := range list {
		resMap[quote.TargetKind] = append(resMap[quote.TargetKind], quote.Target)
	}
	return resMap, nil
}

func (q *quoteStore) GetTargetQuote(ctx context.Context, target int, targetKind quote_entry.QuoteTargetKindType) (map[quote_entry.QuoteKindType][]int, error) {
	list, err := q.ListQuery(ctx, "`target` = ? and `target_kind` = ?", []interface{}{target, targetKind}, "")
	if err != nil {
		return nil, err
	}
	resMap := make(map[quote_entry.QuoteKindType][]int)
	for _, quote := range list {
		resMap[quote.Kind] = append(resMap[quote.Kind], quote.Source)
	}
	return resMap, nil
}

func newQuoteStore(db store.IDB) IQuoteStore {
	quote := &quoteStore{BaseStore: store.CreateStore[quote_entry.Quote](db)}
	return quote
}
