package enum

import "fmt"

type PublishType int

const (
	PublishTypeNone = iota
	PublishTypeUnPublished
	PublishTypePublished
	PublishTypeDefect
	PublishTypeAll
)

var (
	publishNames = map[PublishType]string{
		PublishTypeNone:        "NONE",
		PublishTypeUnPublished: "UNPUBLISHED", //未发布
		PublishTypePublished:   "PUBLISHED",   //已发布
		PublishTypeDefect:      "DEFECT",      //缺失
	}
	PublishIndex = map[string]PublishType{
		"UNPUBLISHED": PublishTypeUnPublished,
		"PUBLISHED":   PublishTypePublished,
		"DEFECT":      PublishTypeDefect,
	}
)

func (e PublishType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprint("\"", e.String(), "\"")), nil
}

func (e PublishType) String() string {
	if e >= PublishTypeAll {
		return "unknown"
	}
	return publishNames[e]
}
