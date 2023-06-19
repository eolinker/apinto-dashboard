//go:build mysql

package main

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/db_migrator"
	"github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/apinto-dashboard/store/store_mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

func initDB() {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", GetDBUserName(), GetDBPassword(), GetDBIp(), GetDBPort(), GetDbName())
	dialector := mysql.Open(dns)
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}
	sqlDb, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDb.SetConnMaxLifetime(time.Second * 9)
	sqlDb.SetMaxOpenConns(200)
	sqlDb.SetMaxIdleConns(200)
	db_migrator.InitSql(db)
	store.InitStoreDB(store_mysql.NewMyDB(db, store.DBInfo{
		Addr: fmt.Sprintf("%s:%d", GetDBIp(), GetDBPort()),
		User: GetDBUserName(),
		DB:   GetDbName(),
	}))

}
