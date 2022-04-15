package dao

type ActivityLogDao interface {
	GetLogList(offset, limit int) ([]*LogEntity, int64, error)
	InsertLog(user, content, operation, object string, args []*Arg) error
}
