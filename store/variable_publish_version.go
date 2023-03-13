package store

import (
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/entry"
)

type IVariablePublishVersionStore interface {
	IBaseStore[entry.VariablePublishVersion]
}

type variablePublishVersionHandler struct {
}

func (s *variablePublishVersionHandler) Kind() string {
	return "publish_variable"
}

func (s *variablePublishVersionHandler) Encode(sv *entry.VariablePublishVersion) *entry.Version {

	v := new(entry.Version)
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

func (s *variablePublishVersionHandler) Decode(v *entry.Version) *entry.VariablePublishVersion {
	sv := new(entry.VariablePublishVersion)
	sv.Id = v.Id
	sv.ClusterId = v.Target
	sv.Operator = v.Operator
	sv.NamespaceId = v.NamespaceID
	sv.CreateTime = v.CreateTime
	_ = json.Unmarshal(v.Data, &sv.VariablePublishVersionConfig)

	return sv
}

func newVariablePublishVersionStore(db IDB) IVariablePublishVersionStore {
	var h BaseKindHandler[entry.VariablePublishVersion, entry.Version] = &variablePublishVersionHandler{}
	return CreateBaseKindStore(h, db)
}
