package dynamic_store

import (
	dynamic_entry "github.com/eolinker/apinto-dashboard/modules/dynamic/dynamic-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IDynamicQuoteStore interface {
	store.IBaseStore[dynamic_entry.DynamicQuote]
}

type dynamicQuoteStore struct {
	*store.BaseStore[dynamic_entry.DynamicQuote]
}

func newDynamicQuoteStore(db store.IDB) IDynamicQuoteStore {
	return &dynamicQuoteStore{BaseStore: store.CreateStore[dynamic_entry.DynamicQuote](db)}
}
