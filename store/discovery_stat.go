package store

import (
	"github.com/eolinker/apinto-dashboard/entry"
)

type IDiscoveryStatStore interface {
	IBaseStore[entry.DiscoveryStat]
}

type DiscoveryStatKindHandler struct {
}

func (s *DiscoveryStatKindHandler) Kind() string {
	return "discovery"
}

func (s *DiscoveryStatKindHandler) Encode(sv *entry.DiscoveryStat) *entry.Stat {
	stat := new(entry.Stat)

	stat.Tag = sv.DiscoveryID
	stat.Kind = s.Kind()
	stat.Version = sv.VersionID

	return stat
}

func (s *DiscoveryStatKindHandler) Decode(stat *entry.Stat) *entry.DiscoveryStat {
	ds := new(entry.DiscoveryStat)

	ds.DiscoveryID = stat.Tag
	ds.VersionID = stat.Version

	return ds
}

func newDiscoveryStatStore(db IDB) IDiscoveryStatStore {
	var h BaseKindHandler[entry.DiscoveryStat, entry.Stat] = new(DiscoveryStatKindHandler)
	return CreateBaseKindStore(h, db)
}
