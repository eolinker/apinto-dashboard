package store

import (
	"github.com/eolinker/apinto-dashboard/entry/notice-entry"
	"github.com/eolinker/apinto-dashboard/entry/stat-entry"
)

type INoticeChannelStatStore interface {
	IBaseStore[notice_entry.NoticeChannelStat]
}

type NoticeChannelKindHandler struct {
}

func (s *NoticeChannelKindHandler) Kind() string {
	return "notice_channel"
}

func (s *NoticeChannelKindHandler) Encode(sv *notice_entry.NoticeChannelStat) *stat_entry.Stat {
	stat := new(stat_entry.Stat)

	stat.Tag = sv.NoticeChannelID
	stat.Kind = s.Kind()
	stat.Version = sv.VersionID

	return stat
}

func (s *NoticeChannelKindHandler) Decode(stat *stat_entry.Stat) *notice_entry.NoticeChannelStat {
	ds := new(notice_entry.NoticeChannelStat)

	ds.NoticeChannelID = stat.Tag
	ds.VersionID = stat.Version

	return ds
}

func newNoticeChannelStatStore(db IDB) INoticeChannelStatStore {
	var h BaseKindHandler[notice_entry.NoticeChannelStat, stat_entry.Stat] = new(NoticeChannelKindHandler)
	return CreateBaseKindStore(h, db)
}
