package notice_service

import (
	driver_manager "github.com/eolinker/apinto-dashboard/driver"
	"github.com/eolinker/apinto-dashboard/modules/notice"
)

type driverNoticeChannel struct {
	*driver_manager.DriverManager[notice.IDriverNoticeChannel]
}

func newNoticeChannelDriverManager() notice.INoticeChannelDriverManager {
	return &driverNoticeChannel{DriverManager: driver_manager.CreateDriverManager[notice.IDriverNoticeChannel]()}
}

// NewNoticeChannelDriverManager
// todo mock 这个文件会报错 Loading input failed: don't know how to mock method of type *ast.IndexExpr 可能是泛型不支持
// 单独提供个方法处理
func NewNoticeChannelDriverManager() notice.INoticeChannelDriverManager {
	return &driverNoticeChannel{DriverManager: driver_manager.CreateDriverManager[notice.IDriverNoticeChannel]()}
}
