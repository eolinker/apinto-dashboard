package enum

import "fmt"

const (
	HeaderOptTypeAdd    = "ADD"    //新增或修改
	HeaderOptTypeDelete = "DELETE" //删除

	MatchPositionHeader = "header"
	MatchPositionQuery  = "query"
	MatchPositionCookie = "cookie"

	MatchTypeEqual   = "EQUAL"   //全等匹配
	MatchTypePrefix  = "PREFIX"  //前缀匹配
	MatchTypeSuffix  = "SUFFIX"  //后缀匹配
	MatchTypeSubstr  = "SUBSTR"  //子串匹配
	MatchTypeUneuqal = "UNEQUAL" //非等匹配
	MatchTypeNull    = "NULL"    //空值匹配
	MatchTypeExist   = "EXIST"   //存在匹配
	MatchTypeUnexist = "UNEXIST" //不存在匹配
	MatchTypeRegexp  = "REGEXP"  //区分大小写的正则匹配
	MatchTypeRegexpG = "REGEXPG" //不区分大小写的匹配
	MatchTypeAny     = "ANY"     //任意匹配

	MethodGET     = "GET"
	MethodPOST    = "POST"
	MethodPUT     = "PUT"
	MethodDELETE  = "DELETE"
	MethodPATCH   = "PATCH"
	MethodHEAD    = "HEAD"
	MethodOPTIONS = "OPTIONS"

	RestfulLabel = "{rest}"

	//来源类型
	SourceSelfBuild = "self-build" //自建
	SourceImport    = "import"     //导入
	SourceSync      = "sync"       //同步
)

var (
	sourceTypeTitles = map[string]string{
		SourceSelfBuild: "自建",
		SourceImport:    "导入",
		SourceSync:      "同步",
	}
)

func GetSourceTitle(sourceType string) string {
	return sourceTypeTitles[sourceType]
}

// ImportStatusType 导入或同步API的状态， 正常或冲突
type ImportStatusType int

const (
	ImportStatusTypeNone = iota
	ImportStatusTypeNormal
	ImportStatusTypeConflict
	ImportStatusTypeInvalidPath
	ImportStatusTypeALL
)

var (
	importStatusTypeNames = map[ImportStatusType]string{
		ImportStatusTypeNone:        "NONE",
		ImportStatusTypeNormal:      "normal",
		ImportStatusTypeConflict:    "conflict",
		ImportStatusTypeInvalidPath: "invalidPath",
	}
)

func (l ImportStatusType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprint("\"", l.String(), "\"")), nil
}

func (l ImportStatusType) String() string {
	if l >= ImportStatusTypeALL {
		return "unknown"
	}
	return importStatusTypeNames[l]
}
