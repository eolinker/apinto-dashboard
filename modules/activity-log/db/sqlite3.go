package db

import (
	"database/sql"
	"github.com/eolinker/apinto-dashboard/modules/activity-log/db/module"
	_ "github.com/mattn/go-sqlite3"
)

func InitDB() error {
	db, err := sql.Open("sqlite3", "./sqlite.db")
	if err != nil {
		return err
	}
	err = module.NewLogModule(db)
	return err
}
