package entry

import "time"

// AuditLog 集群配置信息表
type AuditLog struct {
	Id          int       `gorm:"type:int(11);size:11;not null;auto_increment;primary_key;column:id;comment:主键ID"`
	NamespaceId int       `gorm:"type:int(11);size:11;not null;column:namespace;comment:工作空间" json:"namespace,omitempty"`
	UserID      int       `gorm:"type:int(11);size:11;not null;column:user_id;comment:用户ID"`
	Username    string    `gorm:"size:20;not null;column:username;comment:用户名"`
	IP          string    `gorm:"size:20;not null;column:ip;comment:ip地址"`
	OperateType int       `gorm:"type:tinyint(1);size:1;not null;column:operate;comment:操作类型 1.创建 2.编辑 3.删除 4.发布"`
	Kind        string    `gorm:"size:20;not null;column:kind;comment:操作对象"`
	Object      string    `gorm:"type:text;not null;column:object;comment:对象信息"`
	URL         string    `gorm:"type:text;not null;column:url;comment:请求url,包括query参数"`
	StartTime   time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;column:start_time;comment:请求开始时间"`
	EndTime     time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;column:end_time;comment:请求结束时间"`
	UserAgent   string    `gorm:"type:text;not null;column:user_agent;comment:user-agent"`
	Body        string    `gorm:"type:text;not null;column:body;comment:请求内容"`
	Err         string    `gorm:"type:text;not null;column:error;comment:错误信息"`
}

func (*AuditLog) TableName() string {
	return "audit_log"
}

func (a *AuditLog) IdValue() int {
	return a.Id
}
