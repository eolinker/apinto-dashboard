package notice_service

import "github.com/eolinker/eosc/common/bean"

func init() {
	iNoticeChannelService := newNoticeChannelService()

	bean.Injection(&iNoticeChannelService)
}
