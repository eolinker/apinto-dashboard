package dynamic_store

import (
	"encoding/json"

	dynamic_entry "github.com/eolinker/apinto-dashboard/modules/dynamic/dynamic-entry"

	publish_entry "github.com/eolinker/apinto-dashboard/modules/base/publish-entry"

	"github.com/eolinker/apinto-dashboard/store"
)

type IDynamicPublishHistoryStore interface {
	store.BasePublishHistoryStore[dynamic_entry.DynamicPublishHistory]
}

type dynamicPublishHistoryHandler struct {
	kind string
}

func (s *dynamicPublishHistoryHandler) Kind() string {
	return s.kind
}

func (s *dynamicPublishHistoryHandler) Encode(sr *dynamic_entry.DynamicPublishHistory) *publish_entry.PublishHistory {

	val, _ := json.Marshal(sr.Publish)
	history := &publish_entry.PublishHistory{
		Id:          sr.Id,
		Kind:        s.Kind(),
		ClusterId:   sr.ClusterId,
		NamespaceId: sr.NamespaceId,
		Target:      sr.Target,
		VersionId:   sr.VersionId,
		Data:        string(val),
		Desc:        sr.Desc,
		OptType:     sr.OptType,
		OptTime:     sr.CreateTime,
		VersionName: sr.VersionName,
		Operator:    sr.Operator,
	}
	return history
}

func (s *dynamicPublishHistoryHandler) Decode(r *publish_entry.PublishHistory) *dynamic_entry.DynamicPublishHistory {
	var cfg dynamic_entry.DynamicPublishConfig
	json.Unmarshal([]byte(r.Data), &cfg)
	history := &dynamic_entry.DynamicPublishHistory{
		Id:          r.Id,
		VersionName: r.VersionName,
		Desc:        r.Desc,
		NamespaceId: r.NamespaceId,
		ClusterId:   r.ClusterId,
		VersionId:   r.VersionId,
		Publish:     &cfg,
		CreateTime:  r.OptTime,
		OptType:     r.OptType,
		Operator:    r.Operator,
	}
	return history
}

func NewDynamicPublishHistoryStore(db store.IDB, kind string) IDynamicPublishHistoryStore {
	var historyHandler store.BaseKindHandler[dynamic_entry.DynamicPublishHistory, publish_entry.PublishHistory] = &dynamicPublishHistoryHandler{
		kind: kind,
	}
	return store.CreatePublishHistory[dynamic_entry.DynamicPublishHistory](historyHandler, db)
}
