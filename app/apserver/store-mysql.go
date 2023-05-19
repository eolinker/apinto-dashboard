//go:build mysql

package main

import (
	"context"
	"fmt"
	"github.com/eolinker/apinto-dashboard/db_migrator"
	_ "github.com/eolinker/apinto-dashboard/initialize"
	"github.com/eolinker/apinto-dashboard/store"
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
	store.InitStoreDB(&myDB{
		db: db,
		info: store.DBInfo{
			Addr: fmt.Sprintf("%s:%s", GetDBIp(), GetDBPort()),
			User: GetDBUserName(),
			DB:   GetDbName(),
		},
	})

}

var (
	_ store.IDB = (*myDB)(nil)
)

type myDB struct {
	db   *gorm.DB
	info store.DBInfo
}

var txContextKey = _TxContextKey{}

type _TxContextKey struct {
}

func (m *myDB) Info() store.DBInfo {
	return m.info
}

func (m *myDB) DB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txContextKey).(*gorm.DB); ok {
		return tx
	}
	return m.db.WithContext(ctx)
}
func (m *myDB) IsTxCtx(ctx context.Context) bool {
	if _, ok := ctx.Value(txContextKey).(*gorm.DB); ok {
		return ok
	}
	return false
}
