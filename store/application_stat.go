package store

import (
	"github.com/eolinker/apinto-dashboard/entry"
)

type IApplicationStatStore interface {
	IBaseStore[entry.ApplicationStat]
}

type applicationStatKindHandler struct {
}

func (s *applicationStatKindHandler) Kind() string {
	return "application"
}

func (s *applicationStatKindHandler) Encode(sv *entry.ApplicationStat) *entry.Stat {
	stat := new(entry.Stat)

	stat.Tag = sv.ApplicationID
	stat.Kind = s.Kind()
	stat.Version = sv.VersionID

	return stat
}

func (s *applicationStatKindHandler) Decode(stat *entry.Stat) *entry.ApplicationStat {
	ds := new(entry.ApplicationStat)

	ds.ApplicationID = stat.Tag
	ds.VersionID = stat.Version

	return ds
}

func newApplicationStatStore(db IDB) IApplicationStatStore {
	var h BaseKindHandler[entry.ApplicationStat, entry.Stat] = new(applicationStatKindHandler)
	return CreateBaseKindStore(h, db)
}
