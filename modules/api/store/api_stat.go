package api_store

import (
	"github.com/eolinker/apinto-dashboard/entry"
	api_entry "github.com/eolinker/apinto-dashboard/modules/api/api-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IAPIStatStore interface {
	store.IBaseStore[api_entry.APIStat]
}

type apiHandler struct {
}

func (a *apiHandler) Kind() string {
	return "api"
}

func (a *apiHandler) Encode(as *api_entry.APIStat) *entry.Stat {
	stat := new(entry.Stat)

	stat.Tag = as.APIID
	stat.Kind = a.Kind()
	stat.Version = as.VersionID

	return stat
}

func (a *apiHandler) Decode(stat *entry.Stat) *api_entry.APIStat {
	ds := new(api_entry.APIStat)

	ds.APIID = stat.Tag
	ds.VersionID = stat.Version

	return ds
}

func NewAPIStatStore(db store.IDB) IAPIStatStore {
	var h store.BaseKindHandler[api_entry.APIStat, entry.Stat] = new(apiHandler)
	return store.CreateBaseKindStore(h, db)
}
