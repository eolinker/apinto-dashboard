package entry

import "time"

type WarnHistory struct {
	Id            int       `gorm:"type:int(11);size:11;not null;auto_increment;primary_key;column:id;comment:主键ID"`
	NamespaceID   int       `gorm:"type:int(11);size:11;not null;column:namespace;comment:工作空间"`
	PartitionId   int       `gorm:"type:int(11);size:11;not null;column:partition_id;comment:分区ID"`
	StrategyTitle string    `gorm:"column:strategy_title"`
	ErrMsg        string    `gorm:"column:err_msg"`
	Status        int       `gorm:"column:status"` //发送状态 0未发送 1已发送 2发送失败
	Dimension     string    `gorm:"column:dimension"`
	Quota         string    `gorm:"column:quota"`
	Target        string    `gorm:"column:target"`
	Content       string    `gorm:"column:content"`
	CreateTime    time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;column:create_time;comment:创建时间"`
}

func (WarnHistory) TableName() string {
	return "warn_history"
}

func (n *WarnHistory) IdValue() int {
	return n.Id
}
