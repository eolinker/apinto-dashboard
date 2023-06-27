package namespace_model

import (
	"time"
)

type Namespace struct {
	Id         int
	Name       string
	CreateTime time.Time
}
