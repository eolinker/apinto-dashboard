package store

import (
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/entry"
)

type IVariablePublishHistoryStore interface {
	BasePublishHistoryStore[entry.VariablePublishHistory]
}

type variablePublishHistoryHandler struct {
}

func (s *variablePublishHistoryHandler) Kind() string {
	return "variable"
}

func (s *variablePublishHistoryHandler) Encode(sr *entry.VariablePublishHistory) *entry.PublishHistory {
	val, _ := json.Marshal(sr.VariablePublishHistoryInfo)
	history := &entry.PublishHistory{
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

func (s *variablePublishHistoryHandler) Decode(r *entry.PublishHistory) *entry.VariablePublishHistory {
	val := new(entry.VariablePublishHistoryInfo)
	_ = json.Unmarshal([]byte(r.Data), val)
	history := &entry.VariablePublishHistory{
		Id:                         r.Id,
		VersionName:                r.VersionName,
		Desc:                       r.Desc,
		NamespaceId:                r.NamespaceId,
		ClusterId:                  r.ClusterId,
		VersionId:                  r.VersionId,
		OptTime:                    r.OptTime,
		OptType:                    r.OptType,
		VariablePublishHistoryInfo: *val,
		Operator:                   r.Operator,
	}
	return history
}

func newVariablePublishHistoryStore(db IDB) IVariablePublishHistoryStore {
	var historyHandler BaseKindHandler[entry.VariablePublishHistory, entry.PublishHistory] = new(variablePublishHistoryHandler)
	return createPublishHistory[entry.VariablePublishHistory](historyHandler, db)
}
