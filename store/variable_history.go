package store

import (
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/entry/history-entry"
	"github.com/eolinker/apinto-dashboard/entry/variable-entry"
)

type IVariableHistoryStore interface {
	BaseHistoryStore[variable_entry.VariableHistory]
}

type variableHistoryHandler struct {
}

func (s *variableHistoryHandler) Kind() string {
	return "variable"
}

func (s *variableHistoryHandler) Decode(r *history_entry.History) *variable_entry.VariableHistory {
	oldValue := new(variable_entry.VariableValue)
	_ = json.Unmarshal([]byte(r.OldValue), oldValue)
	newValue := new(variable_entry.VariableValue)
	_ = json.Unmarshal([]byte(r.NewValue), newValue)
	history := &variable_entry.VariableHistory{
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
	var historyHandler DecodeHistory[variable_entry.VariableHistory] = new(variableHistoryHandler)
	return CreateHistory(historyHandler, db, history_entry.HistoryKindVariableGlobal)
}
