package notice_entry

import "time"

type NoticeChannel struct {
	Id          int       `gorm:"type:int(11);size:11;not null;auto_increment;primary_key;column:id;comment:主键ID"`
	NamespaceID int       `gorm:"type:int(11);size:11;not null;column:namespace;dbUniqueIndex:namespace_name;uniqueIndex:namespace_name;comment:工作空间"`
	Name        string    `gorm:"size:36;not null;column:name;dbUniqueIndex:namespace_name;uniqueIndex:namespace_name;comment:名称"`
	Title       string    `gorm:"size:255;not null;column:title;comment:标题"`
	Type        int       `gorm:"type:int(11);size:11;column:type;comment:1.webhook 2.email"`
	Operator    int       `gorm:"type:int(11);size:11;column:operator;comment:更新人/操作人"`
	CreateTime  time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;column:create_time;comment:创建时间"`
	UpdateTime  time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;column:update_time;comment:修改时间"`
}

func (NoticeChannel) TableName() string {
	return "notice_channel"
}

func (n *NoticeChannel) IdValue() int {
	return n.Id
}

type NoticeChannelStat struct {
	NoticeChannelID int
	VersionID       int
}

type NoticeChannelVersion struct {
	Id              int
	NoticeChannelID int
	NamespaceID     int
	NoticeChannelVersionConfig
	Operator   int
	CreateTime time.Time
}

type NoticeChannelVersionConfig struct {
	Config string `json:"config"`
}

func (d *NoticeChannelVersion) SetVersionId(id int) {
	d.Id = id
}
