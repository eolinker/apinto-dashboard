package entry

import "time"

// CommonGroup 通用目录表
type CommonGroup struct {
	Id          int       `gorm:"type:int(11);size:11;not null;auto_increment;primary_key;column:id;comment:主键ID"`
	Uuid        string    `gorm:"size:36;default:null;column:uuid;dbUniqueIndex:uuid_index;uniqueIndex:uuid_index;comment:uuid"`
	NamespaceId int       `gorm:"type:int(11);size:11;not null;index:type;column:namespace;comment:工作空间"`
	Type        string    `gorm:"size:50;not null;index:type;column:type;comment:类型"`
	TagID       int       `gorm:"type:int(11);size:11;not null;index:type;column:tag;comment:根据type区分是哪个表的ID"`
	Name        string    `gorm:"size:50;not null;column:name;comment:名称"`
	ParentId    int       `gorm:"type:int(11);size:11;not null;default:0;column:parent_id;comment:父级目录ID"`
	Operator    int       `gorm:"type:int(11);size:11;default:null;column:operator;comment:更新人/操作人"`
	Sort        int       `gorm:"type:int(11);size:11;not null;column:sort;comment:排序"`
	CreateTime  time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;column:create_time;comment:创建时间"`
	UpdateTime  time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;column:update_time;comment:修改时间"`
}

func (c *CommonGroup) IdValue() int {
	return c.Id
}

func (CommonGroup) TableName() string {
	return "common_group"
}
