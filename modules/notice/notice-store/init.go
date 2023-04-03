package notice_store

import (
	"github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/eosc/common/bean"
)

func init() {
	store.RegisterStore(func(db store.IDB) {
		iNoticeChannelStore := newNoticeChannelStore(db)
		iNoticeChannelStatStore := newNoticeChannelStatStore(db)
		iNoticeChannelVersionStore := newNoticeChannelVersionStore(db)
		bean.Injection(&iNoticeChannelStore)
		bean.Injection(&iNoticeChannelStatStore)
		bean.Injection(&iNoticeChannelVersionStore)
	})
}
