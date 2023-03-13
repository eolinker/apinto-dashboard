package enum

import "fmt"

type OnlineStatus int

const (
	OnlineStatusNone = iota
	OnlineStatusNOTGOONLINE
	OnlineStatusOFFLINE
	OnlineStatusGOONLINE
	OnlineStatusTOUPDATE
	OnlineStatusAll
)

var (
	onlineStatusNames = map[OnlineStatus]string{
		OnlineStatusNone:        "NONE",
		OnlineStatusGOONLINE:    "GOONLINE",    //上线
		OnlineStatusOFFLINE:     "OFFLINE",     //下线
		OnlineStatusNOTGOONLINE: "NOTGOONLINE", //未上线
		OnlineStatusTOUPDATE:    "TOUPDATE",    //待更新
	}

	OnlineStatusIndex = map[string]OnlineStatus{
		"NONE":        OnlineStatusNone,
		"GOONLINE":    OnlineStatusGOONLINE,    //上线
		"OFFLINE":     OnlineStatusOFFLINE,     //下线
		"NOTGOONLINE": OnlineStatusNOTGOONLINE, //未上线
		"TOUPDATE":    OnlineStatusTOUPDATE,    //待更新
	}
)

func (e OnlineStatus) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprint("\"", e.String(), "\"")), nil
}

func (e OnlineStatus) String() string {
	if e >= OnlineStatusAll {
		return "unknown"
	}
	return onlineStatusNames[e]
}
