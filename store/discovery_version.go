package store

import (
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/entry"
)

type IDiscoveryVersionStore interface {
	IBaseStore[entry.DiscoveryVersion]
}

type discoveryVersionKindHandler struct {
}

func (s *discoveryVersionKindHandler) Kind() string {
	return "discovery"
}

func (s *discoveryVersionKindHandler) Encode(dv *entry.DiscoveryVersion) *entry.Version {

	v := new(entry.Version)
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

func (s *discoveryVersionKindHandler) Decode(v *entry.Version) *entry.DiscoveryVersion {
	sv := new(entry.DiscoveryVersion)
	sv.Id = v.Id
	sv.DiscoveryID = v.Target
	sv.Operator = v.Operator
	sv.NamespaceID = v.NamespaceID
	sv.CreateTime = v.CreateTime
	_ = json.Unmarshal(v.Data, &sv.DiscoveryVersionConfig)

	return sv
}

func newDiscoveryVersionStore(db IDB) IDiscoveryVersionStore {
	var h BaseKindHandler[entry.DiscoveryVersion, entry.Version] = new(discoveryVersionKindHandler)
	return CreateBaseKindStore(h, db)
}
