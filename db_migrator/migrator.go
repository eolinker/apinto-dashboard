package db_migrator

import (
	"errors"
	"fmt"
	"github.com/eolinker/apinto-dashboard/entry"
	"gorm.io/gorm"
	"strings"
	"time"
)

func InitSql(db *gorm.DB) {

	ddlList, dmlList := GetSqlList()

	hasTable := db.Migrator().HasTable(&entry.FlywaySchemaHistory{})
	if !hasTable {
		err := db.Exec(entry.FlywaySchemaHistory{}.CreateTableSql()).Error
		if err != nil {
			panic(err)
		}
	}

	rollbackColumnSql := make([]string, 0) //需要回滚修改字段的sql //倒叙执行
	rollbackTableSql := make([]string, 0)  //需要回滚创表的sql //倒叙执行
	isRollback := false

	for _, info := range ddlList {
		history := new(entry.FlywaySchemaHistory)
		db.Where("`version_num` = ? and `type` = 'ddl'", info.VersionNum).First(history)
		if history.Success {
			continue
		}

		history.VersionNum = info.VersionNum
		history.Version = info.Version
		history.Script = info.Script
		history.Md5 = info.Md5
		history.Type = "ddl"
		history.CreateTime = time.Now()

		err := db.Save(history).Error
		if err != nil {
			panic(err)
		}

		for _, data := range info.SqlDataList {
			if err = db.Exec(data.Sql).Error; err != nil {
				isRollback = true
				goto rollback
			}

			sql, _ := getRollbackColumnSql(db, data.TableName, data.Sql)
			if sql != "" {
				rollbackColumnSql = append(rollbackColumnSql, sql)
			} else {
				if isCreateTableSql(data.Sql) {
					rollbackTableSql = append(rollbackTableSql, fmt.Sprintf("DROP TABLE IF EXISTS `%s`", data.TableName))
				}
			}
		}

		history.Success = true
		db.Save(history)
	}

	for _, info := range dmlList {

		history := new(entry.FlywaySchemaHistory)
		db.Where("`version_num` = ? and `type` = 'dml'", info.VersionNum).First(history)

		if history.Md5 != "" && history.Md5 != info.Md5 {
			panic(fmt.Sprintf("%s文件遭到修改", info.Script))
		}

		if history.Success {
			continue
		}

		history.VersionNum = info.VersionNum
		history.Version = info.Version
		history.Script = info.Script
		history.Type = "dml"
		history.Md5 = info.Md5
		history.CreateTime = time.Now()
		db.Save(history)

		err := db.Transaction(func(tx *gorm.DB) error {

			for _, data := range info.SqlDataList {
				if err := tx.Exec(data.Sql).Error; err != nil {
					return err
				}
			}

			return nil
		})
		if err != nil {
			panic(err)
		}
		history.Success = true
		db.Save(history)
	}

rollback:
	if isRollback {
		for i := len(rollbackColumnSql) - 1; i >= 0; i-- {
			db.Exec(rollbackColumnSql[i])
		}
		for i := len(rollbackTableSql) - 1; i >= 0; i-- {
			db.Exec(rollbackTableSql[i])
		}
	}
}

func getRollbackColumnSql(db *gorm.DB, tableName, sql string) (string, error) {

	split := strings.Split(sql, "`")

	upper := strings.ToUpper(split[2])
	isOptColumn := strings.Count(upper, "COLUMN") > 0

	opt := ""
	if isOptColumn {
		if strings.Count(upper, "ADD") > 0 {
			opt = "ADD"
		} else if strings.Count(upper, "MODIFY") > 0 {
			opt = "MODIFY"
		} else if strings.Count(upper, "DROP") > 0 {
			opt = "DROP"
		}

	}
	if opt == "" {
		return "", nil
	}
	return getColumnSql(db, tableName, split[3], opt)
}

func isCreateTableSql(sql string) bool {
	split := strings.Split(sql, "`")
	upper := strings.ToUpper(split[0])
	return strings.Count(upper, "CREATE") > 0 && strings.Count(upper, "TABLE") > 0
}

func getColumnSql(db *gorm.DB, tableName, columnName, opt string) (string, error) {
	columns, err := db.Migrator().ColumnTypes(tableName)
	if err != nil {
		return "", err
	}

	if opt == "ADD" {
		//操作如果是add 说明之前是没有这个字段的 可直接返回DROP 命令
		return fmt.Sprintf("ALTER TABLE `%s` DROP  COLUMN`%s`;", tableName, columnName), nil
	}

	for _, column := range columns {
		if column.Name() == columnName {
			columnType, _ := column.ColumnType()
			isNull := "NOT NULL"
			nullable, _ := column.Nullable()
			if nullable {
				isNull = "NULL"
			}

			defaultValue, ok := column.DefaultValue()
			if ok {
				defaultValue = "DEFAULT " + defaultValue
			}

			comment, commentOk := column.Comment()
			if commentOk {
				comment = fmt.Sprintf("COMMENT '%s'", comment)
			}

			columnSql := ""
			if opt == "DROP" {
				columnSql = fmt.Sprintf("ALTER TABLE `%s` ADD  COLUMN `%s` %s %s %s %s;", tableName, column.Name(), columnType, isNull, defaultValue, comment)
			} else {
				columnSql = fmt.Sprintf("ALTER TABLE `%s` MODIFY COLUMN `%s` %s %s %s %s;", tableName, column.Name(), columnType, isNull, defaultValue, comment)
			}
			return columnSql, nil
		}
	}

	return "", errors.New("找不到字段")
}
