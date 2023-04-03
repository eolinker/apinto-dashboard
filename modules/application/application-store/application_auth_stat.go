package application_store

import (
	"github.com/eolinker/apinto-dashboard/modules/application/application-entry"
	"github.com/eolinker/apinto-dashboard/modules/base/stat-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IApplicationAuthStatStore interface {
	store.IBaseStore[application_entry.ApplicationAuthStat]
}

type applicationAuthStatKindHandler struct {
}

func (s *applicationAuthStatKindHandler) Kind() string {
	return "application_auth"
}

func (s *applicationAuthStatKindHandler) Encode(sv *application_entry.ApplicationAuthStat) *stat_entry.Stat {
	stat := new(stat_entry.Stat)

	stat.Tag = sv.ApplicationAuthId
	stat.Kind = s.Kind()
	stat.Version = sv.VersionID

	return stat
}

func (s *applicationAuthStatKindHandler) Decode(stat *stat_entry.Stat) *application_entry.ApplicationAuthStat {
	ds := new(application_entry.ApplicationAuthStat)

	ds.ApplicationAuthId = stat.Tag
	ds.VersionID = stat.Version

	return ds
}

func newApplicationAuthStatStore(db store.IDB) IApplicationAuthStatStore {
	var h store.BaseKindHandler[application_entry.ApplicationAuthStat, stat_entry.Stat] = new(applicationAuthStatKindHandler)
	return store.CreateBaseKindStore(h, db)
}
