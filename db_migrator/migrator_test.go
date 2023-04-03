package db_migrator

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"testing"
	"time"
)

func getDb() *gorm.DB {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", "root", "123456", "192.168.198.133", 3306, "apinto")

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
	return db
}
func TestGetMigrator(t *testing.T) {
	type args struct {
	}
	var tests = []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitSql(getDb())

		})
	}
}
