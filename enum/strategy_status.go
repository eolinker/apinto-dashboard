package enum

import "fmt"

type StrategyOnlineStatus int

const (
	StrategyOnlineStatusNone        = iota
	StrategyOnlineStatusTOUPDATE    //待更新
	StrategyOnlineStatusGOONLINE    //已上线
	StrategyOnlineStatusTODELETE    //待删除
	StrategyOnlineStatusNOTGOONLINE //未上线
	StrategyOnlineStatusAll
)

var (
	StrategyOnlineStatusNames = map[StrategyOnlineStatus]string{
		StrategyOnlineStatusNone:        "NONE",
		StrategyOnlineStatusGOONLINE:    "GOONLINE",    //上线
		StrategyOnlineStatusTODELETE:    "TODELETE",    //下线
		StrategyOnlineStatusNOTGOONLINE: "NOTGOONLINE", //未上线
		StrategyOnlineStatusTOUPDATE:    "TOUPDATE",    //待更新
	}

	StrategyOnlineStatusIndex = map[string]StrategyOnlineStatus{}
)

func init() {
	for i := 0; i < StrategyOnlineStatusAll; i++ {
		e := StrategyOnlineStatus(i)
		StrategyOnlineStatusIndex[e.String()] = e
	}
}

func (e StrategyOnlineStatus) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprint("\"", e.String(), "\"")), nil
}

func (e StrategyOnlineStatus) String() string {
	if e >= StrategyOnlineStatusAll {
		return "unknown"
	}
	return StrategyOnlineStatusNames[e]
}
