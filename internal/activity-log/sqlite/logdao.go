package sqlite

import (
	"database/sql"
	"encoding/json"
	"fmt"
	apinto "github.com/eolinker/apinto-dashboard"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
	"time"
)

type activityLogDao struct {
	db *sql.DB
}

func (a *activityLogDao) GetLogList(offset, limit int, user, operation, target string, startUnix, endUnix int64) ([]*apinto.LogEntity, int64, error) {
	db := a.db

	list := make([]*apinto.LogEntity, 0, limit)
	var totalNum int64

	//拼接sql语句
	totalSQL := "select count(id) from `activityLog` where timestamp >= ?"
	listSQL := "select `user`,`ip`,`operation`,`target`,`content`,`args`,`timestamp` from `activityLog` where timestamp >= ?"
	params := make([]interface{}, 0, 2)
	params = append(params, startUnix)

	if endUnix != 0 {
		//因为得到的时间戳精度为年月日，因此结束时间需要改为当天的最后一秒
		endUnix = endUnix + 86400 - 1
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

	if target != "" {
		totalSQL = totalSQL + " and target = ?"
		listSQL = listSQL + " and target = ?"
		params = append(params, target)
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
			ip, content, argsJson string
			timestamp             int64
		)

		err = rows.Scan(&user, &ip, &operation, &target, &content, &argsJson, &timestamp)
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
			IP:        ip,
			Operation: operation,
			Target:    target,
			Content:   content,
			Args:      args,
		}

		list = append(list, entity)
	}

	return list, totalNum, nil
}

func (a *activityLogDao) Add(user, ip, operation, target, content string, args []*apinto.Arg) error {
	db := a.db

	timestamp := time.Now().Unix()
	details, _ := json.Marshal(args)
	if len(details) == 0 {
		details = []byte{'[', ']'}
	}
	const sqlStatement = "INSERT INTO `activityLog` (`user`,`ip`,`operation`,`target`,`content`,`args`,`timestamp`) VALUES (?,?,?,?,?,?,?);"
	_, err := db.Exec(sqlStatement, user, ip, operation, target, content, details, timestamp)
	return err
}

func (a *activityLogDao) initTable() error {
	db := a.db
	const sqlStatement = "CREATE TABLE IF NOT EXISTS `activityLog` (\n `id` INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,\n `user` VARCHAR(20),\n `ip` VARCHAR(20),\n `operation` VARCHAR(20),\n `target` VARCHAR(20),\n `content` VARCHAR(255),\n `args` TEXT,\n `timestamp` INTEGER NOT NULL\n );\n CREATE INDEX IF NOT EXISTS `index_timestamp` ON `activityLog` (`timestamp`);\n"
	_, err := db.Exec(sqlStatement)

	return err
}

func NewActivityDao(filePath string) (ISqliteHandler, error) {
	err := createDir(filePath)
	if err != nil {
		return nil, fmt.Errorf("Create DB Dir Fail: %s ", err.Error())
	}
	db, err := sql.Open("sqlite", filePath)
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

func createDir(filePath string) error {
	dir := filepath.Dir(filePath)
	if apinto.IsDirExist(dir) {
		return nil
	}
	return os.Mkdir(dir, os.ModeDir)
}
