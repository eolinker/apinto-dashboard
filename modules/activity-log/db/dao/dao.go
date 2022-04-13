package dao

type ActivityLogDao interface {
	GetLogList(page, pageSize int) []*LogEntity
	InsertLog(user, content string, args []*Arg) error
}
