package dao

import (
	"database/sql"
	"fmt"
)

type activityLogDao struct {
	db *sql.DB
}

func (a *activityLogDao) GetLogList(page, pageSize int) ([]*LogEntity, int64) {
	panic("implement me")
}

func (a *activityLogDao) InsertLog(user, content, operation, object string, args []*Arg) error {
	panic("implement me")
}

func (a *activityLogDao) initTable() error {
	db := a.db
	const sqlStatement = "CREATE TABLE IF NOT EXISTS `activityLog` (\n `id` INTEGER PRIMARY KEY AUTOINCREMENT,\n `user` VARCHAR(20),\n `operation` VARCHAR(20),\n `object` VARCHAR(20),\n `content` VARCHAR(255),\n `args` TEXT,\n `timestamp` INTEGER NOT NULL\n );\n CREATE INDEX `index_timestamp` ON `activityLog` (`timestamp`);\n"
	_, err := db.Exec(sqlStatement)
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
