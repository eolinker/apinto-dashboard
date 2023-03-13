package store

import (
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/entry"
)

type IApplicationVersionStore interface {
	IBaseStore[entry.ApplicationVersion]
}

type applicationVersionKindHandler struct {
}

func (s *applicationVersionKindHandler) Kind() string {
	return "application"
}

func (s *applicationVersionKindHandler) Encode(sv *entry.ApplicationVersion) *entry.Version {

	v := new(entry.Version)
	v.Id = sv.Id
	v.Kind = s.Kind()
	v.Target = sv.ApplicationID
	v.Operator = sv.Operator
	v.NamespaceID = sv.NamespaceID
	v.CreateTime = sv.CreateTime
	bytes, _ := json.Marshal(sv.ApplicationVersionConfig)
	v.Data = bytes

	return v
}

func (s *applicationVersionKindHandler) Decode(v *entry.Version) *entry.ApplicationVersion {
	sv := new(entry.ApplicationVersion)
	sv.Id = v.Id
	sv.ApplicationID = v.Target
	sv.Operator = v.Operator
	sv.NamespaceID = v.NamespaceID
	sv.CreateTime = v.CreateTime
	_ = json.Unmarshal(v.Data, &sv.ApplicationVersionConfig)

	return sv
}

func newApplicationVersionStore(db IDB) IApplicationVersionStore {
	var h BaseKindHandler[entry.ApplicationVersion, entry.Version] = new(applicationVersionKindHandler)
	return CreateBaseKindStore(h, db)
}
