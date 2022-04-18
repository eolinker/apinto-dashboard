package dao

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type activityLogDao struct {
	db *sql.DB
}

func (a *activityLogDao) GetLogList(offset, limit int) ([]*LogEntity, int64, error) {
	list := make([]*LogEntity, 0, limit)

	db := a.db

	var totalNum int64

	//查询符合要求的总行数
	err := db.QueryRow("select count(id) from `activityLog`").Scan(&totalNum)
	if err != nil {
		if err == sql.ErrNoRows {
			return list, 0, nil
		}
		return list, 0, fmt.Errorf("GetLogList.GetTotalRows Fail. %s", err)
	}

	rows, err := db.Query("select `user`,`operation`,`object`,`content`,`args`,`timestamp` from `activityLog` ORDER BY `timestamp` DESC limit ? offset ?", limit, offset)
	if err != nil {
		return list, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			user, operation, object, content, argsJson string
			timestamp                                  int64
		)

		err = rows.Scan(&user, &operation, &object, &content, &argsJson, &timestamp)
		if err != nil {
			return list, 0, err
		}

		args := make([]Arg, 0)
		err = json.Unmarshal([]byte(argsJson), &args)
		if err != nil {
			return list, 0, fmt.Errorf("GetLogList.Unmarshal args Fail. %s", err)
		}

		entity := &LogEntity{
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

func (a *activityLogDao) InsertLog(user, operation, object, content string, args []*Arg) error {
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

	for i := 1; i <= 10; i++ {
		err = a.InsertLog("admin", "创建", fmt.Sprintf("demoRouter_%d", i), "创建demoRouter", []*Arg{{Key: "avc", Value: "123"}, {Key: "zzz", Value: 321}})
		if err != nil {
			return err
		}
		err = a.InsertLog("admin", "删除", fmt.Sprintf("demoRouter_%d", i), "删除", nil)
		if err != nil {
			return err
		}
	}
	return err
}

func NewActivityDao(db *sql.DB) (ActivityLogDao, error) {
	a := &activityLogDao{
		db: db,
	}

	err := a.initTable()
	if err != nil {
		return nil, fmt.Errorf("activityLogDao initTable Fail. %s", err)
	}
	return a, nil
}
