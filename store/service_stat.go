package store

import (
	"github.com/eolinker/apinto-dashboard/entry"
)

type IServiceStatStore interface {
	IBaseStore[entry.ServiceStat]
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

func newServiceStatStore(db IDB) IServiceStatStore {
	var h BaseKindHandler[entry.ServiceStat, entry.Stat] = new(ServiceStatKindHandler)
	return CreateBaseKindStore(h, db)
}
