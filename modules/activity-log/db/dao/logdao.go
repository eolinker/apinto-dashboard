package dao

import (
	"database/sql"
	"fmt"
)

type activityLogDao struct {
	db *sql.DB
}

func (a *activityLogDao) GetLogList(page, pageSize int) []*LogEntity {
	panic("implement me")
}

func (a *activityLogDao) InsertLog(user, content string, args []*Arg) error {
	panic("implement me")
}

func (a *activityLogDao) initTable() error {
	//panic("implement me")
	return nil
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
