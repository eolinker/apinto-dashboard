package api_store

import (
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/modules/api/api-entry"
	"github.com/eolinker/apinto-dashboard/modules/base/publish-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IApiPublishHistoryStore interface {
	store.BasePublishHistoryStore[api_entry.ApiPublishHistory]
}

type apiPublishHistoryHandler struct {
}

func (s *apiPublishHistoryHandler) Kind() string {
	return "api"
}

func (s *apiPublishHistoryHandler) Encode(sr *api_entry.ApiPublishHistory) *publish_entry.PublishHistory {
	val, _ := json.Marshal(sr.APIVersionConfig)
	history := &publish_entry.PublishHistory{
		Id:          sr.Id,
		Kind:        s.Kind(),
		ClusterId:   sr.ClusterId,
		NamespaceId: sr.NamespaceId,
		Target:      sr.ClusterId,
		VersionId:   sr.VersionId,
		Data:        string(val),
		Desc:        sr.Desc,
		OptType:     sr.OptType,
		OptTime:     sr.OptTime,
		VersionName: sr.VersionName,
		Operator:    sr.Operator,
	}
	return history
}

func (s *apiPublishHistoryHandler) Decode(r *publish_entry.PublishHistory) *api_entry.ApiPublishHistory {
	val := new(api_entry.APIVersionConfig)
	_ = json.Unmarshal([]byte(r.Data), val)
	history := &api_entry.ApiPublishHistory{
		Id:               r.Id,
		VersionName:      r.VersionName,
		Desc:             r.Desc,
		NamespaceId:      r.NamespaceId,
		ClusterId:        r.ClusterId,
		VersionId:        r.VersionId,
		OptTime:          r.OptTime,
		OptType:          r.OptType,
		APIVersionConfig: *val,
		Operator:         r.Operator,
	}
	return history
}

func newApiPublishHistoryStore(db store.IDB) IApiPublishHistoryStore {
	var historyHandler store.BaseKindHandler[api_entry.ApiPublishHistory, publish_entry.PublishHistory] = new(apiPublishHistoryHandler)
	return store.CreatePublishHistory[api_entry.ApiPublishHistory](historyHandler, db)
}
