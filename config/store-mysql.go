//go:build mysql

package config

import (
	"fmt"

	store "github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/apinto-dashboard/store/store_mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	slog "log"
	"os"
	"time"
)

func InitDb() {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", getDBUserName(), getDBPassword(), getDBIp(), getDBPort(), getDBName())
	dialector := mysql.Open(dns)
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.New(slog.New(os.Stdout, "\r\n", slog.LstdFlags), logger.Config{
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
		Addr: fmt.Sprintf("%s:%d", getDBIp(), getDBPort()),
		User: getDBUserName(),
		DB:   getDBName(),
	}))
}
