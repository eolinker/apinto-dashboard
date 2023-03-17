//go:build mysql
// +build mysql

package main

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/store"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

func init() {
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
	store.InitStoreDB(db)

}
