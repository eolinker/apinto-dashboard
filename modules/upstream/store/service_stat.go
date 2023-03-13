package upstream_store

import (
	"github.com/eolinker/apinto-dashboard/entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IServiceStatStore interface {
	store.IBaseStore[entry.ServiceStat]
}

type ServiceStatKindHandler struct {
}

func (s *ServiceStatKindHandler) Kind() string {
	return "service"
}

func (s *ServiceStatKindHandler) Encode(sv *entry.ServiceStat) *entry.Stat {
	stat := new(entry.Stat)

	stat.Tag = sv.ServiceId
	stat.Kind = "service"
	stat.Version = sv.VersionId

	return stat
}

func (s *ServiceStatKindHandler) Decode(stat *entry.Stat) *entry.ServiceStat {
	ss := new(entry.ServiceStat)

	ss.ServiceId = stat.Tag
	ss.VersionId = stat.Version

	return ss
}

func newServiceStatStore(db store.IDB) IServiceStatStore {
	var h store.BaseKindHandler[entry.ServiceStat, entry.Stat] = new(ServiceStatKindHandler)
	return store.CreateBaseKindStore(h, db)
}
