package model

import (
	"github.com/eolinker/apinto-dashboard/entry"
)

type Cluster struct {
	*entry.Cluster
	Status int //1正常 2部分正常 3异常
}
