package sqlite

import (
	"database/sql"
	"encoding/json"
	"fmt"
	apinto "github.com/eolinker/apinto-dashboard"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type activityLogDao struct {
	db *sql.DB
}

func (a *activityLogDao) GetLogList(offset, limit int, user, operation, object string, startUnix, endUnix int64) ([]*apinto.LogEntity, int64, error) {
	db := a.db

	list := make([]*apinto.LogEntity, 0, limit)
	var totalNum int64

	//拼接sql语句
	totalSQL := "select count(id) from `activityLog` where timestamp >= ?"
	listSQL := "select `user`,`operation`,`object`,`content`,`args`,`timestamp` from `activityLog` where timestamp >= ?"
	params := make([]interface{}, 0, 2)
	params = append(params, startUnix)

	if endUnix != 0 {
		totalSQL = totalSQL + " and timestamp <= ?"
		listSQL = listSQL + " and timestamp <= ?"
		params = append(params, endUnix)
	}

	if user != "" {
		totalSQL = totalSQL + " and user = ?"
		listSQL = listSQL + " and user = ?"
		params = append(params, user)
	}

	if operation != "" {
		totalSQL = totalSQL + " and operation = ?"
		listSQL = listSQL + " and operation = ?"
		params = append(params, operation)
	}

	if object != "" {
		totalSQL = totalSQL + " and object = ?"
		listSQL = listSQL + " and object = ?"
		params = append(params, object)
	}

	//查询符合要求的总行数
	err := db.QueryRow(totalSQL, params...).Scan(&totalNum)
	if err != nil {
		if err == sql.ErrNoRows {
			return list, 0, nil
		}
		return list, 0, fmt.Errorf("GetLogList.GetTotalRows Fail. %s", err)
	}

	listSQL = listSQL + "ORDER BY `timestamp` DESC limit ? offset ?"
	params = append(params, limit, offset)

	//查询列表
	rows, err := db.Query(listSQL, params...)
	if err != nil {
		return list, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			content, argsJson string
			timestamp         int64
		)

		err = rows.Scan(&user, &operation, &object, &content, &argsJson, &timestamp)
		if err != nil {
			return list, 0, err
		}

		args := make([]*apinto.Arg, 0)
		err = json.Unmarshal([]byte(argsJson), &args)
		if err != nil {
			return list, 0, fmt.Errorf("GetLogList.Unmarshal args Fail. %s", err)
		}

		entity := &apinto.LogEntity{
			Time:      time.Unix(timestamp, 0).Format("2006-01-02 15:04:05"),
			User:      user,
			Operation: operation,
			Object:    object,
			Content:   content,
			Args:      args,
		}

		list = append(list, entity)
	}

	return list, totalNum, nil
}

func (a *activityLogDao) Add(user, operation, object, content string, args []*apinto.Arg) error {
	db := a.db

	timestamp := time.Now().Unix()
	details, _ := json.Marshal(args)
	if len(details) == 0 {
		details = []byte{'[', ']'}
	}
	const sqlStatement = "INSERT INTO `activityLog` (`user`,`operation`,`object`,`content`,`args`,`timestamp`) VALUES (?,?,?,?,?,?);"
	_, err := db.Exec(sqlStatement, user, operation, object, content, details, timestamp)
	return err
}

func (a *activityLogDao) initTable() error {
	db := a.db
	const sqlStatement = "CREATE TABLE IF NOT EXISTS `activityLog` (\n `id` INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,\n `user` VARCHAR(20),\n `operation` VARCHAR(20),\n `object` VARCHAR(20),\n `content` VARCHAR(255),\n `args` TEXT,\n `timestamp` INTEGER NOT NULL\n );\n CREATE INDEX IF NOT EXISTS `index_timestamp` ON `activityLog` (`timestamp`);\n"
	_, err := db.Exec(sqlStatement)

	//刷数据
	//var time int64 = 1650420291
	//for i := 1; i <= 14; i++ {
	//	err = a.InsertLog("admin", "登录", "", "admin成功登录", nil, time)
	//	if err != nil {
	//		return err
	//	}
	//	err = a.InsertLog("admin", "创建", fmt.Sprintf("demoRouter_%d", i), "创建demoRouter", []*Arg{{Key: "avc", Value: "123"}, {Key: "zzz", Value: 321}, {Key: "object", Value: map[string]interface{}{"a": 1, "b": "2"}}}, time+5)
	//	if err != nil {
	//		return err
	//	}
	//	err = a.InsertLog("admin", "删除", fmt.Sprintf("demoRouter_%d", i), "删除", nil, time+10)
	//	if err != nil {
	//		return err
	//	}
	//	err = a.InsertLog("admin", "登录", "", "admin登出", nil, time+15)
	//	if err != nil {
	//		return err
	//	}
	//
	//	err = a.InsertLog("eolink", "登录", "", "eolink成功登录", nil, time+20)
	//	if err != nil {
	//		return err
	//	}
	//	err = a.InsertLog("eolink", "创建", fmt.Sprintf("demoRouter_%d", i), "创建demoRouter", []*Arg{{Key: "avc", Value: "123"}, {Key: "zzz", Value: 321}, {Key: "object", Value: map[string]interface{}{"a": 1, "b": "2"}}}, time+25)
	//	if err != nil {
	//		return err
	//	}
	//	err = a.InsertLog("eolink", "删除", fmt.Sprintf("demoRouter_%d", i), "删除", nil, time+30)
	//	if err != nil {
	//		return err
	//	}
	//	err = a.InsertLog("eolink", "登录", "", "eolink登出", nil, time+35)
	//	if err != nil {
	//		return err
	//	}
	//	time += 86400
	//}
	return err
}

func NewActivityDao(file string) (ISqliteHandler, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}
	a := &activityLogDao{
		db: db,
	}

	err = a.initTable()
	if err != nil {
		return nil, fmt.Errorf("activityLogDao initTable Fail. %s", err)
	}
	return a, nil
}
