package application_store

import (
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/entry/application-entry"
	"github.com/eolinker/apinto-dashboard/entry/version-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IApplicationVersionStore interface {
	store.IBaseStore[application_entry.ApplicationVersion]
}

type applicationVersionKindHandler struct {
}

func (s *applicationVersionKindHandler) Kind() string {
	return "application"
}

func (s *applicationVersionKindHandler) Encode(sv *application_entry.ApplicationVersion) *version_entry.Version {

	v := new(version_entry.Version)
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

func (s *applicationVersionKindHandler) Decode(v *version_entry.Version) *application_entry.ApplicationVersion {
	sv := new(application_entry.ApplicationVersion)
	sv.Id = v.Id
	sv.ApplicationID = v.Target
	sv.Operator = v.Operator
	sv.NamespaceID = v.NamespaceID
	sv.CreateTime = v.CreateTime
	_ = json.Unmarshal(v.Data, &sv.ApplicationVersionConfig)

	return sv
}

func newApplicationVersionStore(db store.IDB) IApplicationVersionStore {
	var h store.BaseKindHandler[application_entry.ApplicationVersion, version_entry.Version] = new(applicationVersionKindHandler)
	return store.CreateBaseKindStore(h, db)
}
