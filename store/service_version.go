package store

import (
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/entry"
)

type IServiceVersionStore interface {
	IBaseStore[entry.ServiceVersion]
}

type serviceVersionKindHandler struct {
}

func (s *serviceVersionKindHandler) Kind() string {
	return "service"
}

func (s *serviceVersionKindHandler) Encode(sv *entry.ServiceVersion) *entry.Version {

	v := new(entry.Version)
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

func (s *serviceVersionKindHandler) Decode(v *entry.Version) *entry.ServiceVersion {
	sv := new(entry.ServiceVersion)
	sv.Id = v.Id
	sv.ServiceId = v.Target
	sv.Operator = v.Operator
	sv.NamespaceID = v.NamespaceID
	sv.CreateTime = v.CreateTime
	_ = json.Unmarshal(v.Data, &sv.ServiceVersionConfig)

	return sv
}

func newServiceVersionStore(db IDB) IServiceVersionStore {
	var h BaseKindHandler[entry.ServiceVersion, entry.Version] = new(serviceVersionKindHandler)
	return CreateBaseKindStore(h, db)
}
