package role_store

import (
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/entry/history-entry"
	"github.com/eolinker/apinto-dashboard/entry/role-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IRoleAccessLogStore interface {
	store.BaseHistoryStore[role_entry.RoleAccessLog]
}

type roleAccessLogHandler struct {
}

func (s *roleAccessLogHandler) Kind() string {
	return "role_access"
}

func (s *roleAccessLogHandler) Encode(sr *role_entry.RoleAccessLog) *history_entry.History {
	oldValue, _ := json.Marshal(sr.OldValue)
	newValue, _ := json.Marshal(sr.NewValue)
	log := &history_entry.History{
		Kind:     s.Kind(),
		TargetID: sr.RoleID,
		OldValue: string(oldValue),
		NewValue: string(newValue),
		OptType:  history_entry.OptType(sr.OptType),
		Operator: sr.Operator,
		OptTime:  sr.OptTime,
	}
	return log
}

func (s *roleAccessLogHandler) Decode(r *history_entry.History) *role_entry.RoleAccessLog {
	oldValue := new(role_entry.AccessListLog)
	_ = json.Unmarshal([]byte(r.OldValue), oldValue)
	newValue := new(role_entry.AccessListLog)
	_ = json.Unmarshal([]byte(r.NewValue), newValue)
	history := &role_entry.RoleAccessLog{
		Id:       r.Id,
		Operator: r.Operator,
		RoleID:   r.TargetID,
		OldValue: *oldValue,
		NewValue: *newValue,
		OptType:  int(r.OptType),
		OptTime:  r.OptTime,
	}
	return history
}

func newRoleAccessLogStore(db store.IDB) IRoleAccessLogStore {
	var historyHandler store.DecodeHistory[role_entry.RoleAccessLog] = new(roleAccessLogHandler)
	return store.CreateHistory(historyHandler, db, history_entry.HistoryKindRole)
}
