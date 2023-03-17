package notice_store

import (
	"github.com/eolinker/apinto-dashboard/modules/base/stat-entry"
	"github.com/eolinker/apinto-dashboard/modules/notice/notice-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type INoticeChannelStatStore interface {
	store.IBaseStore[notice_entry.NoticeChannelStat]
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

func newNoticeChannelStatStore(db store.IDB) INoticeChannelStatStore {
	var h store.BaseKindHandler[notice_entry.NoticeChannelStat, stat_entry.Stat] = new(NoticeChannelKindHandler)
	return store.CreateBaseKindStore(h, db)
}
