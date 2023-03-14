package upstream_store

import (
	"github.com/eolinker/apinto-dashboard/entry/stat-entry"
	"github.com/eolinker/apinto-dashboard/entry/upstream-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IServiceStatStore interface {
	store.IBaseStore[upstream_entry.ServiceStat]
}

type ServiceStatKindHandler struct {
}

func (s *ServiceStatKindHandler) Kind() string {
	return "service"
}

func (s *ServiceStatKindHandler) Encode(sv *upstream_entry.ServiceStat) *stat_entry.Stat {
	stat := new(stat_entry.Stat)

	stat.Tag = sv.ServiceId
	stat.Kind = "service"
	stat.Version = sv.VersionId

	return stat
}

func (s *ServiceStatKindHandler) Decode(stat *stat_entry.Stat) *upstream_entry.ServiceStat {
	ss := new(upstream_entry.ServiceStat)

	ss.ServiceId = stat.Tag
	ss.VersionId = stat.Version

	return ss
}

func newServiceStatStore(db store.IDB) IServiceStatStore {
	var h store.BaseKindHandler[upstream_entry.ServiceStat, stat_entry.Stat] = new(ServiceStatKindHandler)
	return store.CreateBaseKindStore(h, db)
}
