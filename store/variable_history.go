package store

import (
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/entry"
)

type IVariableHistoryStore interface {
	BaseHistoryStore[entry.VariableHistory]
}

type variableHistoryHandler struct {
}

func (s *variableHistoryHandler) Kind() string {
	return "variable"
}

func (s *variableHistoryHandler) Decode(r *entry.History) *entry.VariableHistory {
	oldValue := new(entry.VariableValue)
	_ = json.Unmarshal([]byte(r.OldValue), oldValue)
	newValue := new(entry.VariableValue)
	_ = json.Unmarshal([]byte(r.NewValue), newValue)
	history := &entry.VariableHistory{
		Id:          r.Id,
		NamespaceId: r.NamespaceID,
		VariableId:  r.TargetID,
		OptTime:     r.OptTime,
		OptType:     r.OptType,
		OldValue:    *oldValue,
		NewValue:    *newValue,
		Operator:    r.Operator,
	}
	return history
}

func newVariableHistoryStore(db IDB) IVariableHistoryStore {
	var historyHandler DecodeHistory[entry.VariableHistory] = new(variableHistoryHandler)
	return CreateHistory(historyHandler, db, entry.HistoryKindVariableGlobal)
}
