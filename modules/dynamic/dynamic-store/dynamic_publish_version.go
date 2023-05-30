package dynamic_store

import (
	"encoding/json"

	dynamic_entry "github.com/eolinker/apinto-dashboard/modules/dynamic/dynamic-entry"

	"github.com/eolinker/apinto-dashboard/modules/base/version-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IDynamicPublishVersionStore interface {
	store.IBaseStore[dynamic_entry.DynamicPublishVersion]
}

type dynamicPublishVersionHandler struct {
	kind string
}

func (s *dynamicPublishVersionHandler) Kind() string {
	return s.kind
}

func (s *dynamicPublishVersionHandler) Encode(sv *dynamic_entry.DynamicPublishVersion) *version_entry.Version {
	bytes, _ := json.Marshal(sv.Publish)
	v := new(version_entry.Version)
	v.Id = sv.Id
	v.Kind = s.Kind()
	v.Target = sv.ClusterId
	v.Operator = sv.Operator
	v.CreateTime = sv.CreateTime
	v.NamespaceID = sv.NamespaceId
	v.Data = bytes

	return v
}

func (s *dynamicPublishVersionHandler) Decode(v *version_entry.Version) *dynamic_entry.DynamicPublishVersion {
	var cfg dynamic_entry.DynamicPublishConfig
	json.Unmarshal(v.Data, &cfg)
	sv := new(dynamic_entry.DynamicPublishVersion)
	sv.Id = v.Id
	sv.ClusterId = v.Target
	sv.Operator = v.Operator
	sv.NamespaceId = v.NamespaceID
	sv.CreateTime = v.CreateTime
	sv.Publish = &cfg
	return sv
}

func NewDynamicPublishVersionStore(db store.IDB, kind string) IDynamicPublishVersionStore {
	var h store.BaseKindHandler[dynamic_entry.DynamicPublishVersion, version_entry.Version] = &dynamicPublishVersionHandler{
		kind: kind,
	}
	return store.CreateBaseKindStore(h, db)
}
