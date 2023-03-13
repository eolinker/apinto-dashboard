package store

import (
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/entry"
)

type INoticeChannelVersionStore interface {
	IBaseStore[entry.NoticeChannelVersion]
}

type noticeChannelVersionKindHandler struct {
}

func (s *noticeChannelVersionKindHandler) Kind() string {
	return "notice_channel"
}

func (s *noticeChannelVersionKindHandler) Encode(dv *entry.NoticeChannelVersion) *entry.Version {

	v := new(entry.Version)
	v.Id = dv.Id
	v.Kind = s.Kind()
	v.NamespaceID = dv.NamespaceID
	v.Target = dv.NoticeChannelID
	v.Operator = dv.Operator
	v.CreateTime = dv.CreateTime
	bytes, _ := json.Marshal(dv.NoticeChannelVersionConfig)
	v.Data = bytes

	return v
}

func (s *noticeChannelVersionKindHandler) Decode(v *entry.Version) *entry.NoticeChannelVersion {
	sv := new(entry.NoticeChannelVersion)
	sv.Id = v.Id
	sv.NoticeChannelID = v.Target
	sv.Operator = v.Operator
	sv.NamespaceID = v.NamespaceID
	sv.CreateTime = v.CreateTime
	_ = json.Unmarshal(v.Data, &sv.NoticeChannelVersionConfig)

	return sv
}

func newNoticeChannelVersionStore(db IDB) INoticeChannelVersionStore {
	var h BaseKindHandler[entry.NoticeChannelVersion, entry.Version] = new(noticeChannelVersionKindHandler)
	return CreateBaseKindStore(h, db)
}
