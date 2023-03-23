package model

import "time"

type WarnHistoryInfo struct {
	StrategyTitle string
	Dimension     string
	Quota         string
	Target        string
	Content       string
	ErrMsg        string
	Status        int //发送状态 0未发送 1已发送 2发送失败  3部分成功
	CreateTime    time.Time
}
