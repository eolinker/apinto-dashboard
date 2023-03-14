package store

import (
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/entry/publish-entry"
	"github.com/eolinker/apinto-dashboard/entry/variable-entry"
)

type IVariablePublishHistoryStore interface {
	BasePublishHistoryStore[variable_entry.VariablePublishHistory]
}

type variablePublishHistoryHandler struct {
}

func (s *variablePublishHistoryHandler) Kind() string {
	return "variable"
}

func (s *variablePublishHistoryHandler) Encode(sr *variable_entry.VariablePublishHistory) *publish_entry.PublishHistory {
	val, _ := json.Marshal(sr.VariablePublishHistoryInfo)
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

func (s *variablePublishHistoryHandler) Decode(r *publish_entry.PublishHistory) *variable_entry.VariablePublishHistory {
	val := new(variable_entry.VariablePublishHistoryInfo)
	_ = json.Unmarshal([]byte(r.Data), val)
	history := &variable_entry.VariablePublishHistory{
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
	var historyHandler BaseKindHandler[variable_entry.VariablePublishHistory, publish_entry.PublishHistory] = new(variablePublishHistoryHandler)
	return createPublishHistory[variable_entry.VariablePublishHistory](historyHandler, db)
}
