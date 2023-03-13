package store

import (
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/entry"
)

type IRoleAccessLogStore interface {
	BaseHistoryStore[entry.RoleAccessLog]
}

type roleAccessLogHandler struct {
}

func (s *roleAccessLogHandler) Kind() string {
	return "role_access"
}

func (s *roleAccessLogHandler) Encode(sr *entry.RoleAccessLog) *entry.History {
	oldValue, _ := json.Marshal(sr.OldValue)
	newValue, _ := json.Marshal(sr.NewValue)
	log := &entry.History{
		Kind:     s.Kind(),
		TargetID: sr.RoleID,
		OldValue: string(oldValue),
		NewValue: string(newValue),
		OptType:  entry.OptType(sr.OptType),
		Operator: sr.Operator,
		OptTime:  sr.OptTime,
	}
	return log
}

func (s *roleAccessLogHandler) Decode(r *entry.History) *entry.RoleAccessLog {
	oldValue := new(entry.AccessListLog)
	_ = json.Unmarshal([]byte(r.OldValue), oldValue)
	newValue := new(entry.AccessListLog)
	_ = json.Unmarshal([]byte(r.NewValue), newValue)
	history := &entry.RoleAccessLog{
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

func newRoleAccessLogStore(db IDB) IRoleAccessLogStore {
	var historyHandler DecodeHistory[entry.RoleAccessLog] = new(roleAccessLogHandler)
	return createHistory(historyHandler, db, entry.HistoryKindRole)
}
