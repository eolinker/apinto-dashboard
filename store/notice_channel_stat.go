package store

import "github.com/eolinker/apinto-dashboard/entry"

type INoticeChannelStatStore interface {
	IBaseStore[entry.NoticeChannelStat]
}

type NoticeChannelKindHandler struct {
}

func (s *NoticeChannelKindHandler) Kind() string {
	return "notice_channel"
}

func (s *NoticeChannelKindHandler) Encode(sv *entry.NoticeChannelStat) *entry.Stat {
	stat := new(entry.Stat)

	stat.Tag = sv.NoticeChannelID
	stat.Kind = s.Kind()
	stat.Version = sv.VersionID

	return stat
}

func (s *NoticeChannelKindHandler) Decode(stat *entry.Stat) *entry.NoticeChannelStat {
	ds := new(entry.NoticeChannelStat)

	ds.NoticeChannelID = stat.Tag
	ds.VersionID = stat.Version

	return ds
}

func newNoticeChannelStatStore(db IDB) INoticeChannelStatStore {
	var h BaseKindHandler[entry.NoticeChannelStat, entry.Stat] = new(NoticeChannelKindHandler)
	return CreateBaseKindStore(h, db)
}
