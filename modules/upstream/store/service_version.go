package upstream_store

import (
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/entry/upstream-entry"
	"github.com/eolinker/apinto-dashboard/entry/version-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IServiceVersionStore interface {
	store.IBaseStore[upstream_entry.ServiceVersion]
}

type serviceVersionKindHandler struct {
}

func (s *serviceVersionKindHandler) Kind() string {
	return "service"
}

func (s *serviceVersionKindHandler) Encode(sv *upstream_entry.ServiceVersion) *version_entry.Version {

	v := new(version_entry.Version)
	v.Id = sv.Id
	v.Kind = s.Kind()
	v.Target = sv.ServiceId
	v.Operator = sv.Operator
	v.NamespaceID = sv.NamespaceID
	v.CreateTime = sv.CreateTime
	bytes, _ := json.Marshal(sv.ServiceVersionConfig)
	v.Data = bytes

	return v
}

func (s *serviceVersionKindHandler) Decode(v *version_entry.Version) *upstream_entry.ServiceVersion {
	sv := new(upstream_entry.ServiceVersion)
	sv.Id = v.Id
	sv.ServiceId = v.Target
	sv.Operator = v.Operator
	sv.NamespaceID = v.NamespaceID
	sv.CreateTime = v.CreateTime
	_ = json.Unmarshal(v.Data, &sv.ServiceVersionConfig)

	return sv
}

func newServiceVersionStore(db store.IDB) IServiceVersionStore {
	var h store.BaseKindHandler[upstream_entry.ServiceVersion, version_entry.Version] = new(serviceVersionKindHandler)
	return store.CreateBaseKindStore(h, db)
}
