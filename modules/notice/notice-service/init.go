package notice_service

import (
	"github.com/eolinker/eosc/common/bean"
)

func init() {
	iNoticeChannelService := newNoticeChannelService()
	noticeChannelDriverManager := newNoticeChannelDriverManager()
	bean.Injection(&iNoticeChannelService)
	bean.Injection(&noticeChannelDriverManager)
	bean.AddInitializingBeanFunc(func() {

		err := iNoticeChannelService.InitChannelDriver()
		if err != nil {
			panic(err)
		}
	})
}
