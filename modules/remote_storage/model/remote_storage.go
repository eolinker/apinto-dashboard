package model

import "time"

type RemoteStorage struct {
	Module   string
	Key      string
	Object   interface{}
	UpdateAt time.Time
	CreateAt time.Time
}
