package warn_entry

import "time"

// WarnStrategy 告警策略表
type WarnStrategy struct {
	Id          int       `gorm:"type:int(11);size:11;not null;auto_increment;primary_key;column:id;comment:主键ID"`
	NamespaceID int       `gorm:"type:int(11);size:11;not null;column:namespace;dbUniqueIndex:namespace_partition_title;uniqueIndex:namespace_partition_title;comment:工作空间"`
	PartitionId int       `gorm:"type:int(11);size:11;not null;column:partition_id;dbUniqueIndex:namespace_partition_title;uniqueIndex:namespace_partition_title;comment:分区ID"`
	Title       string    `gorm:"size:36;not null;column:title;dbUniqueIndex:namespace_partition_title;uniqueIndex:namespace_partition_title;comment:名称"`
	UUID        string    `gorm:"column:uuid;dbUniqueIndex:uuid"`
	Desc        string    `gorm:"column:desc"`
	IsEnable    bool      `gorm:"column:is_enable"` //是否启用
	Dimension   string    `gorm:"column:dimension"` //告警维度 api/service/cluster/partition
	Quota       string    `gorm:"column:quota"`     //告警指标
	Every       int       `gorm:"column:every"`     //统计时间 单位分钟
	Config      string    `gorm:"column:config"`
	Operator    int       `gorm:"type:int(11);size:11;column:operator;comment:更新人/操作人"`
	CreateTime  time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;column:create_time;comment:创建时间"`
	UpdateTime  time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;column:update_time;comment:修改时间"`
}

func (WarnStrategy) TableName() string {
	return "warn_strategy"
}

func (n *WarnStrategy) IdValue() int {
	return n.Id
}
