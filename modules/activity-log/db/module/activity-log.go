package module

import (
	"database/sql"
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/modules/activity-log/db/dao"
)

var (
	activityLogModule dao.ActivityLogDao
)

func NewLogModule(db *sql.DB) (err error) {
	activityLogModule, err = dao.NewActivityDao(db)
	return
}

func GetLogList(page, pageSize int) []byte {
	list := activityLogModule.GetLogList(page, pageSize)

	data, _ := json.Marshal(list)
	return data
}

func InsertLog(user, content string, args []*dao.Arg) error {
	return activityLogModule.InsertLog(user, content, args)
}
