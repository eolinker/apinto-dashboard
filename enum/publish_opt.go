package enum

import "fmt"

func init() {
	for i := 0; i < PublishOptTypeAll; i++ {
		e := PublishOptType(i)
		publishOptTypeIndex[e.String()] = e
	}
}

//PublishOptType 环境变量发布操作类型
type PublishOptType int

const (
	PublishOptTypeNone = iota
	PublishOptTypePublish
	PublishOptTypeRollback
	PublishOptTypeAll
)

var (
	publishOptTypeNames = map[PublishOptType]string{
		PublishOptTypeNone:     "NONE",
		PublishOptTypePublish:  "PUBLISH",
		PublishOptTypeRollback: "ROLLBACK",
	}
	publishOptTypeIndex = map[string]PublishOptType{}
)

func (e PublishOptType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprint("\"", e.String(), "\"")), nil
}

func (e PublishOptType) String() string {
	if e >= PublishOptTypeAll {
		return "unknown"
	}
	return publishOptTypeNames[e]
}
