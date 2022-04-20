package dao

type ActivityLogDao interface {
	GetLogList(offset, limit int, user, operation, object string, startUnix, endUnix int64) ([]*LogEntity, int64, error)
	InsertLog(user, content, operation, object string, args []*Arg) error
}
