package entry

import "gorm.io/gorm"

type RemoteKeyObject struct {
	Module string      `gorm:"uniqueIndex:module_key;index;not null;size:255;column:module;comment:模块名"`
	Key    string      `gorm:"uniqueIndex:module_key;not null;size:32;column:key"` //TODO key要多长
	Object interface{} `gorm:"type:text;serializer:json;column:object"`
	gorm.Model
}

func (c *RemoteKeyObject) TableName() string {
	return "remote_storage"
}

func (c *RemoteKeyObject) IdValue() int {
	return int(c.ID)
}
