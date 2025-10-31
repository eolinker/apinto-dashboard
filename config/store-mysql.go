//go:build mysql

package config

import (
	"fmt"

	slog "log"
	"os"
	"time"

	store "github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/apinto-dashboard/store/store_mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDb() {
	dialector := mysql.Open(getDBNS())
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.New(slog.New(os.Stderr, "\r\n", slog.LstdFlags), logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		}),
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

	store.InitStoreDB(store_mysql.NewMyDB(db, store.DBInfo{
		Addr: fmt.Sprintf("%s:%d", systemConfig.MysqlConfig.Ip, systemConfig.MysqlConfig.Port),
		User: systemConfig.MysqlConfig.UserName,
		DB:   systemConfig.MysqlConfig.Db,
	}))
}
