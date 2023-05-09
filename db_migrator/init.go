package db_migrator

import (
	"embed"
	"github.com/eolinker/apinto-dashboard/common"
	"io/fs"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

//go:embed sql
var sqlDir embed.FS

type SqlData struct {
	TableName string
	Sql       string
}
type SqlInfo struct {
	Version     string
	VersionNum  int64
	Data        string
	Script      string
	Type        string //ddl,dml
	Md5         string
	SqlDataList []*SqlData
}

func GetSqlList() ([]*SqlInfo, []*SqlInfo) {
	sub, err := fs.Sub(sqlDir, "sql")
	if err != nil {
		panic(err)
	}

	readDir, err := sqlDir.ReadDir("sql")
	if err != nil {
		panic(err)
	}

	sqlInfoDDLList := make([]*SqlInfo, 0)
	sqlInfoDMLList := make([]*SqlInfo, 0)
	for _, entry := range readDir {
		if strings.Count(entry.Name(), ".sql") == 0 {
			continue
		}

		bytes, err := fs.ReadFile(sub, entry.Name())
		if err != nil {
			panic(err)
		}
		entry.Name()
		split := strings.Split(entry.Name(), "_")
		if len(split) < 3 {
			continue
		}

		version := split[0]
		versionNum := versionToNum(version)

		if split[1] == "ddl" {
			sqlInfoDDLList = append(sqlInfoDDLList, &SqlInfo{
				Version:    version,
				VersionNum: versionNum,
				Type:       split[1],
				Data:       string(bytes),
				Script:     entry.Name(),
			})
		} else if split[1] == "dml" {
			sqlInfoDMLList = append(sqlInfoDMLList, &SqlInfo{
				Version:    version,
				VersionNum: versionNum,
				Type:       split[1],
				Data:       string(bytes),
				Script:     entry.Name(),
			})
		}

	}

	sort.Slice(sqlInfoDDLList, func(i, j int) bool {
		return sqlInfoDDLList[i].VersionNum < sqlInfoDDLList[j].VersionNum
	})
	sort.Slice(sqlInfoDMLList, func(i, j int) bool {
		return sqlInfoDMLList[i].VersionNum < sqlInfoDMLList[j].VersionNum
	})

	for _, info := range sqlInfoDDLList {
		sqlArr := strings.Split(info.Data, ";")
		sqlAll := ""
		for _, sql := range sqlArr {
			sql = strings.TrimSpace(sql)
			if sql == "" {
				continue
			}
			sqlAll += sql
			split := strings.Split(sql, "`")
			info.SqlDataList = append(info.SqlDataList, &SqlData{
				TableName: split[1],
				Sql:       sql,
			})
		}
		info.Md5 = common.Md5(sqlAll)
	}

	for _, info := range sqlInfoDMLList {
		sqlArr := strings.Split(info.Data, ";")
		sqlAll := ""
		for _, sql := range sqlArr {
			sql = strings.TrimSpace(sql)
			if sql == "" {
				continue
			}
			sqlAll += sql
			split := strings.Split(sql, "`")
			info.SqlDataList = append(info.SqlDataList, &SqlData{
				TableName: split[1],
				Sql:       sql,
			})
		}
		info.Md5 = common.Md5(sqlAll)
	}

	return sqlInfoDDLList, sqlInfoDMLList
}

// versionToNum 正则匹配更准确
func versionToNum(ver string) int64 {
	if ver == "" {
		return 0
	}
	exp := regexp.MustCompile(`\d+\.\d+\.\d+`)
	vs := exp.FindString(ver)
	if vs == "" {
		return 0
	}
	ss := strings.Split(vs, ".")
	ver1, _ := strconv.Atoi(ss[0])
	ver2, _ := strconv.Atoi(ss[1])
	ver3, _ := strconv.Atoi(ss[2])
	return int64(ver1)*1000000 + int64(ver2)*1000 + int64(ver3)
}
