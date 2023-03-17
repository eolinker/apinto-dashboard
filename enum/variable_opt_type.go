package enum

import "fmt"

type VariableOptType int

const (
	VariableOptTypeNone = iota
	VariableOptTypeNew
	VariableOptTypeModify
	VariableOptTypeDelete
	VariableOptTypeAll
)

var (
	variableOptTypeNames = map[VariableOptType]string{
		VariableOptTypeNone:   "NONE",
		VariableOptTypeNew:    "NEW",
		VariableOptTypeModify: "MODIFY",
		VariableOptTypeDelete: "DELETE",
	}
	VariableOptTypeIndex = map[string]VariableOptType{
		"NEW":    VariableOptTypeNew,
		"MODIFY": VariableOptTypeModify,
		"DELETE": VariableOptTypeDelete,
	}
)

func (e VariableOptType) String() string {
	if e >= VariableOptTypeAll {
		return "unknown"
	}
	return variableOptTypeNames[e]
}

func (e VariableOptType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprint("\"", e.String(), "\"")), nil
}
