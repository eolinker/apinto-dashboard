package variable_store

import (
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/entry/variable-entry"
	"github.com/eolinker/apinto-dashboard/entry/version-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IVariablePublishVersionStore interface {
	store.IBaseStore[variable_entry.VariablePublishVersion]
}

type variablePublishVersionHandler struct {
}

func (s *variablePublishVersionHandler) Kind() string {
	return "publish_variable"
}

func (s *variablePublishVersionHandler) Encode(sv *variable_entry.VariablePublishVersion) *version_entry.Version {

	v := new(version_entry.Version)
	v.Id = sv.Id
	v.Kind = s.Kind()
	v.Target = sv.ClusterId
	v.Operator = sv.Operator
	v.CreateTime = sv.CreateTime
	v.NamespaceID = sv.NamespaceId
	bytes, _ := json.Marshal(sv.VariablePublishVersionConfig)
	v.Data = bytes

	return v
}

func (s *variablePublishVersionHandler) Decode(v *version_entry.Version) *variable_entry.VariablePublishVersion {
	sv := new(variable_entry.VariablePublishVersion)
	sv.Id = v.Id
	sv.ClusterId = v.Target
	sv.Operator = v.Operator
	sv.NamespaceId = v.NamespaceID
	sv.CreateTime = v.CreateTime
	_ = json.Unmarshal(v.Data, &sv.VariablePublishVersionConfig)

	return sv
}

func newVariablePublishVersionStore(db store.IDB) IVariablePublishVersionStore {
	var h store.BaseKindHandler[variable_entry.VariablePublishVersion, version_entry.Version] = &variablePublishVersionHandler{}
	return store.CreateBaseKindStore(h, db)
}
