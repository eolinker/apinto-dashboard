package entry

import "time"

type Middleware struct {
	Id          int       `gorm:"column:id"`
	NamespaceId int       `gorm:"column:namespace_id" json:"namespace_id,omitempty"`
	UUID        string    `gorm:"column:uuid" json:"uuid,omitempty"`
	Prefix      string    `gorm:"column:prefix" json:"prefix,omitempty"`
	Middlewares string    `gorm:"column:middlewares" json:"middlewares,omitempty"`
	Operator    int       `gorm:"column:operator"`
	CreateTime  time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime  time.Time `gorm:"column:update_time" json:"update_time"`
}

func (m *Middleware) IdValue() int {
	return m.Id
}

func (m *Middleware) TableName() string {
	return "middleware_group"
}
