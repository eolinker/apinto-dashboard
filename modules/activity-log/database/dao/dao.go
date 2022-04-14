package dao

type ActivityLogDao interface {
	GetLogList(page, pageSize int) ([]*LogEntity, int64)
	InsertLog(user, content, operation, object string, args []*Arg) error
}
