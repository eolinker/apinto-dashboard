package notice

import "github.com/eolinker/apinto-dashboard/driver"

type INoticeChannelDriverManager interface {
	driver.IDriverManager[IDriverNoticeChannel]
}
