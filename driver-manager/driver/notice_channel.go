package driver

type IDriverNoticeChannel interface {
	SendTo(sends []string, title, msg string) error
}
