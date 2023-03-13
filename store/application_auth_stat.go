package store

import (
	"github.com/eolinker/apinto-dashboard/entry"
)

type IApplicationAuthStatStore interface {
	IBaseStore[entry.ApplicationAuthStat]
}

type applicationAuthStatKindHandler struct {
}

func (s *applicationAuthStatKindHandler) Kind() string {
	return "application_auth"
}

func (s *applicationAuthStatKindHandler) Encode(sv *entry.ApplicationAuthStat) *entry.Stat {
	stat := new(entry.Stat)

	stat.Tag = sv.ApplicationAuthId
	stat.Kind = s.Kind()
	stat.Version = sv.VersionID

	return stat
}

func (s *applicationAuthStatKindHandler) Decode(stat *entry.Stat) *entry.ApplicationAuthStat {
	ds := new(entry.ApplicationAuthStat)

	ds.ApplicationAuthId = stat.Tag
	ds.VersionID = stat.Version

	return ds
}

func newApplicationAuthStatStore(db IDB) IApplicationAuthStatStore {
	var h BaseKindHandler[entry.ApplicationAuthStat, entry.Stat] = new(applicationAuthStatKindHandler)
	return CreateBaseKindStore(h, db)
}
