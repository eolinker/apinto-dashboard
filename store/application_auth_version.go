package store

import (
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/entry"
)

type IApplicationAuthVersionStore interface {
	IBaseStore[entry.ApplicationAuthVersion]
}

type applicationAuthVersionKindHandler struct {
}

func (s *applicationAuthVersionKindHandler) Kind() string {
	return "application_auth"
}

func (s *applicationAuthVersionKindHandler) Encode(sv *entry.ApplicationAuthVersion) *entry.Version {

	v := new(entry.Version)
	v.Id = sv.Id
	v.Kind = s.Kind()
	v.Target = sv.ApplicationAuthID
	v.Operator = sv.Operator
	v.NamespaceID = sv.NamespaceID
	v.CreateTime = sv.CreateTime
	bytes, _ := json.Marshal(sv.ApplicationAuthVersionConfig)
	v.Data = bytes

	return v
}

func (s *applicationAuthVersionKindHandler) Decode(v *entry.Version) *entry.ApplicationAuthVersion {
	sv := new(entry.ApplicationAuthVersion)
	sv.Id = v.Id
	sv.ApplicationAuthID = v.Target
	sv.Operator = v.Operator
	sv.NamespaceID = v.NamespaceID
	sv.CreateTime = v.CreateTime
	_ = json.Unmarshal(v.Data, &sv.ApplicationAuthVersionConfig)

	return sv
}

func newApplicationAuthVersionStore(db IDB) IApplicationAuthVersionStore {
	var h BaseKindHandler[entry.ApplicationAuthVersion, entry.Version] = new(applicationAuthVersionKindHandler)
	return CreateBaseKindStore(h, db)
}
