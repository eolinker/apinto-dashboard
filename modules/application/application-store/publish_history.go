package application_store

import (
	"encoding/json"
	application_entry "github.com/eolinker/apinto-dashboard/modules/application/application-entry"
	"github.com/eolinker/apinto-dashboard/modules/base/publish-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IAppPublishHistoryStore interface {
	store.BasePublishHistoryStore[application_entry.AppPublishHistory]
}

type appPublishHistoryHandler struct {
}

func (s *appPublishHistoryHandler) Kind() string {
	return "application"
}

func (s *appPublishHistoryHandler) Encode(sr *application_entry.AppPublishHistory) *publish_entry.PublishHistory {
	val, _ := json.Marshal(sr.ApplicationVersionConfig)
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
		OptTime:     sr.OptTime,
		VersionName: sr.VersionName,
		Operator:    sr.Operator,
	}
	return history
}

func (s *appPublishHistoryHandler) Decode(r *publish_entry.PublishHistory) *application_entry.AppPublishHistory {
	val := new(application_entry.ApplicationVersionConfig)
	_ = json.Unmarshal([]byte(r.Data), val)
	history := &application_entry.AppPublishHistory{
		Id:                       r.Id,
		VersionName:              r.VersionName,
		ClusterId:                r.ClusterId,
		NamespaceId:              r.NamespaceId,
		Desc:                     r.Desc,
		VersionId:                r.VersionId,
		Target:                   r.Target,
		ApplicationVersionConfig: *val,
		OptType:                  r.OptType,
		Operator:                 r.Operator,
		OptTime:                  r.OptTime,
	}
	return history
}

func newAppPublishHistoryStore(db store.IDB) IAppPublishHistoryStore {
	var historyHandler store.BaseKindHandler[application_entry.AppPublishHistory, publish_entry.PublishHistory] = new(appPublishHistoryHandler)
	return store.CreatePublishHistory[application_entry.AppPublishHistory](historyHandler, db)
}
