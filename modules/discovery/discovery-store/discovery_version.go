package discovery_store

import (
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/modules/base/version-entry"
	"github.com/eolinker/apinto-dashboard/modules/discovery/discovery-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IDiscoveryVersionStore interface {
	store.IBaseStore[discovery_entry.DiscoveryVersion]
}

type discoveryVersionKindHandler struct {
}

func (s *discoveryVersionKindHandler) Kind() string {
	return "discovery"
}

func (s *discoveryVersionKindHandler) Encode(dv *discovery_entry.DiscoveryVersion) *version_entry.Version {

	v := new(version_entry.Version)
	v.Id = dv.Id
	v.Kind = s.Kind()
	v.NamespaceID = dv.NamespaceID
	v.Target = dv.DiscoveryID
	v.Operator = dv.Operator
	v.CreateTime = dv.CreateTime
	bytes, _ := json.Marshal(dv.DiscoveryVersionConfig)
	v.Data = bytes

	return v
}

func (s *discoveryVersionKindHandler) Decode(v *version_entry.Version) *discovery_entry.DiscoveryVersion {
	sv := new(discovery_entry.DiscoveryVersion)
	sv.Id = v.Id
	sv.DiscoveryID = v.Target
	sv.Operator = v.Operator
	sv.NamespaceID = v.NamespaceID
	sv.CreateTime = v.CreateTime
	_ = json.Unmarshal(v.Data, &sv.DiscoveryVersionConfig)

	return sv
}

func newDiscoveryVersionStore(db store.IDB) IDiscoveryVersionStore {
	var h store.BaseKindHandler[discovery_entry.DiscoveryVersion, version_entry.Version] = new(discoveryVersionKindHandler)
	return store.CreateBaseKindStore(h, db)
}
