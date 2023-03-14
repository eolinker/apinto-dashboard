package application_store

import (
	"github.com/eolinker/apinto-dashboard/entry/application-entry"
	"github.com/eolinker/apinto-dashboard/entry/stat-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IApplicationStatStore interface {
	store.IBaseStore[application_entry.ApplicationStat]
}

type applicationStatKindHandler struct {
}

func (s *applicationStatKindHandler) Kind() string {
	return "application"
}

func (s *applicationStatKindHandler) Encode(sv *application_entry.ApplicationStat) *stat_entry.Stat {
	stat := new(stat_entry.Stat)

	stat.Tag = sv.ApplicationID
	stat.Kind = s.Kind()
	stat.Version = sv.VersionID

	return stat
}

func (s *applicationStatKindHandler) Decode(stat *stat_entry.Stat) *application_entry.ApplicationStat {
	ds := new(application_entry.ApplicationStat)

	ds.ApplicationID = stat.Tag
	ds.VersionID = stat.Version

	return ds
}

func newApplicationStatStore(db store.IDB) IApplicationStatStore {
	var h store.BaseKindHandler[application_entry.ApplicationStat, stat_entry.Stat] = new(applicationStatKindHandler)
	return store.CreateBaseKindStore(h, db)
}
