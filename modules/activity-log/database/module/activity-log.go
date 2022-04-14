package module

import (
	"database/sql"
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/modules/activity-log/database/dao"
)

var (
	activityLogModule dao.ActivityLogDao
)

func NewLogModule(db *sql.DB) (err error) {
	activityLogModule, err = dao.NewActivityDao(db)
	return
}

func GetLogList(page, pageSize int) ([]byte, error) {
	list, total, err := activityLogModule.GetLogList(page, pageSize)
	if err != nil {
		return nil, err
	}
	m := make(map[string]interface{})
	m["list"] = list
	m["total_num"] = total

	data, _ := json.Marshal(m)
	return data, nil
}

func InsertLog(user, content, operation, objectName string, args []*dao.Arg) error {
	return activityLogModule.InsertLog(user, content, operation, objectName, args)
}
