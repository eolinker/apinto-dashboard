package enum

import "fmt"

type ChangeOptType int

const (
	ChangeOptTypeNone = iota
	ChangeOptTypeNew
	ChangeOptTypeModify
	ChangeOptTypeDelete
	ChangeOptTypeAll
)

var (
	OptTypeNames = map[ChangeOptType]string{
		ChangeOptTypeNone:   "NONE",
		ChangeOptTypeNew:    "NEW",
		ChangeOptTypeModify: "MODIFY",
		ChangeOptTypeDelete: "DELETE",
	}
	OptTypeIndex = map[string]ChangeOptType{
		"NEW":    ChangeOptTypeNew,
		"MODIFY": ChangeOptTypeModify,
		"DELETE": ChangeOptTypeDelete,
	}
)

func (e ChangeOptType) String() string {
	if e >= ChangeOptTypeAll {
		return "unknown"
	}
	return OptTypeNames[e]
}

func (e ChangeOptType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprint("\"", e.String(), "\"")), nil
}
