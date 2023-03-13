package entry

import "time"

// PublishHistory 发布记录表
type PublishHistory struct {
	Id          int       `gorm:"type:int(11);size:11;not null;auto_increment;primary_key;column:id;comment:主键ID"`
	VersionName string    `gorm:"size:50;not null;column:version_name;comment:版本号"`
	VersionId   int       `gorm:"type:int(11);not null;column:version;comment:版本ID"`
	ClusterId   int       `gorm:"type:int(11);not null;index:cluster;column:cluster;comment:集群ID"`
	NamespaceId int       `gorm:"type:int(11);not null;column:namespace;comment:工作空间"`
	Kind        string    `gorm:"size:50;not null;index:kind_target;column:kind"`
	Target      int       `gorm:"type:int(11);not null;index:kind_target;column:target"`
	OptType     int       `gorm:"type:int(11);not null;column:opt_type;comment:1发布 2回滚"` //1发布 2回滚
	Desc        string    `gorm:"size:255;default:null;column:desc;comment:描述"`
	Data        string    `gorm:"type:text;column:data;comment:数据"`
	Operator    int       `gorm:"type:int(11);size:11;not null;column:operator;comment:操作人"`
	OptTime     time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;column:opt_time;comment:创建时间"`
}

func (*PublishHistory) TableName() string {
	return "publish_history"
}
