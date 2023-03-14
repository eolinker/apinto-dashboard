package store

import (
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/entry/notice-entry"
	"github.com/eolinker/apinto-dashboard/entry/version-entry"
)

type INoticeChannelVersionStore interface {
	IBaseStore[notice_entry.NoticeChannelVersion]
}

type noticeChannelVersionKindHandler struct {
}

func (s *noticeChannelVersionKindHandler) Kind() string {
	return "notice_channel"
}

func (s *noticeChannelVersionKindHandler) Encode(dv *notice_entry.NoticeChannelVersion) *version_entry.Version {

	v := new(version_entry.Version)
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

func (s *noticeChannelVersionKindHandler) Decode(v *version_entry.Version) *notice_entry.NoticeChannelVersion {
	sv := new(notice_entry.NoticeChannelVersion)
	sv.Id = v.Id
	sv.NoticeChannelID = v.Target
	sv.Operator = v.Operator
	sv.NamespaceID = v.NamespaceID
	sv.CreateTime = v.CreateTime
	_ = json.Unmarshal(v.Data, &sv.NoticeChannelVersionConfig)

	return sv
}

func newNoticeChannelVersionStore(db IDB) INoticeChannelVersionStore {
	var h BaseKindHandler[notice_entry.NoticeChannelVersion, version_entry.Version] = new(noticeChannelVersionKindHandler)
	return CreateBaseKindStore(h, db)
}
