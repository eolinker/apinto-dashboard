package api_service

import (
	api_store "github.com/eolinker/apinto-dashboard/modules/api/store"
	"github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/eosc/common/bean"
)

func init() {
	api := NewAPIService()
	store.RegisterStore(api_store.InitStoreDB)
	bean.Injection(&api)
}
