package enum

import "fmt"

type VariableUsageStatus int

const (
	VariableUsageStatusNone = iota
	VariableUsageStatusUnused
	VariableUsageStatusInUse
	VariableUsageStatusAll
)

func init() {
	for i := 0; i < VariableUsageStatusAll; i++ {
		e := VariableUsageStatus(i)
		variableUsageStatusIndex[e.String()] = e
	}
}

var (
	variableUsageStatusNames = map[VariableUsageStatus]string{
		VariableUsageStatusNone:   "NONE",
		VariableUsageStatusUnused: "UNUSED", //空闲
		VariableUsageStatusInUse:  "IN_USE", //使用中
	}
	variableUsageStatusIndex = map[string]VariableUsageStatus{}
)

func (e VariableUsageStatus) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprint("\"", e.String(), "\"")), nil
}

func (e VariableUsageStatus) String() string {
	if e >= VariableUsageStatusAll {
		return "unknown"
	}
	return variableUsageStatusNames[e]
}

//CheckVariableUsageStatus 判断使用状态是否合法
func CheckVariableUsageStatus(status string) bool {
	if _, has := variableUsageStatusIndex[status]; has {
		return true
	}
	return false
}

func GetStatusIndexByName(status string) int {
	idx, has := variableUsageStatusIndex[status]
	if !has {
		return 0
	}
	return int(idx)
}
