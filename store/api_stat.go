package store

import "github.com/eolinker/apinto-dashboard/entry"

type IAPIStatStore interface {
	IBaseStore[entry.APIStat]
}

type apiHandler struct {
}

func (a *apiHandler) Kind() string {
	return "api"
}

func (a *apiHandler) Encode(as *entry.APIStat) *entry.Stat {
	stat := new(entry.Stat)

	stat.Tag = as.APIID
	stat.Kind = a.Kind()
	stat.Version = as.VersionID

	return stat
}

func (a *apiHandler) Decode(stat *entry.Stat) *entry.APIStat {
	ds := new(entry.APIStat)

	ds.APIID = stat.Tag
	ds.VersionID = stat.Version

	return ds
}

func newAPIStatStore(db IDB) IAPIStatStore {
	var h BaseKindHandler[entry.APIStat, entry.Stat] = new(apiHandler)
	return CreateBaseKindStore(h, db)
}
