package driver_manager

import "github.com/eolinker/apinto-dashboard/driver-manager/driver"

type INoticeChannelDriverManager interface {
	IDriverManager[driver.IDriverNoticeChannel]
}

type driverNoticeChannel struct {
	*driverManager[driver.IDriverNoticeChannel]
}

func newNoticeChannelDriverManager() INoticeChannelDriverManager {
	return &driverNoticeChannel{driverManager: createDriverManager[driver.IDriverNoticeChannel]()}
}

// NewNoticeChannelDriverManager
// todo mock 这个文件会报错 Loading input failed: don't know how to mock method of type *ast.IndexExpr 可能是泛型不支持
// 单独提供个方法处理
func NewNoticeChannelDriverManager() INoticeChannelDriverManager {
	return &driverNoticeChannel{driverManager: createDriverManager[driver.IDriverNoticeChannel]()}
}
