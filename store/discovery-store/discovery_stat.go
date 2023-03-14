package discovery_store

import (
	"github.com/eolinker/apinto-dashboard/entry/discovery-entry"
	"github.com/eolinker/apinto-dashboard/entry/stat-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IDiscoveryStatStore interface {
	store.IBaseStore[discovery_entry.DiscoveryStat]
}

type DiscoveryStatKindHandler struct {
}

func (s *DiscoveryStatKindHandler) Kind() string {
	return "discovery"
}

func (s *DiscoveryStatKindHandler) Encode(sv *discovery_entry.DiscoveryStat) *stat_entry.Stat {
	stat := new(stat_entry.Stat)

	stat.Tag = sv.DiscoveryID
	stat.Kind = s.Kind()
	stat.Version = sv.VersionID

	return stat
}

func (s *DiscoveryStatKindHandler) Decode(stat *stat_entry.Stat) *discovery_entry.DiscoveryStat {
	ds := new(discovery_entry.DiscoveryStat)

	ds.DiscoveryID = stat.Tag
	ds.VersionID = stat.Version

	return ds
}

func newDiscoveryStatStore(db store.IDB) IDiscoveryStatStore {
	var h store.BaseKindHandler[discovery_entry.DiscoveryStat, stat_entry.Stat] = new(DiscoveryStatKindHandler)
	return store.CreateBaseKindStore(h, db)
}
