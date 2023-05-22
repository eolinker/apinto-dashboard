package application_entry

import (
	"github.com/eolinker/apinto-dashboard/modules/base/history-entry"
	"time"
)

type ApplicationAuth struct {
	Id            int       `gorm:"type:int(11);size:11;not null;auto_increment;primary_key;column:id;comment:主键ID" json:"id,omitempty"`
	Uuid          string    `gorm:"size:36;not null;column:uuid;uniqueIndex:uuid;comment:uuid" json:"uuid,omitempty"`
	Namespace     int       `gorm:"type:int(11);size:11;not null;column:namespace;comment:工作空间" json:"namespace,omitempty"`
	Application   int       `gorm:"type:int(11);size:11;not null;index:application;column:application;comment:应用ID" json:"application,omitempty"`
	IsTransparent bool      `gorm:"type:tinyint(1);size:1;not null;default:0;column:is_transparent;comment:隐藏鉴权" json:"is_transparent,omitempty"`
	Driver        string    `gorm:"size:30;not null;column:driver;comment:鉴权类型,basic,apikey,aksk,jwt" json:"driver,omitempty"`
	Position      string    `gorm:"size:30;not null;column:position;comment:header,query" json:"position,omitempty"`
	TokenName     string    `gorm:"size:255;not null;column:token_name;comment:tokenName" json:"token_name,omitempty"`
	ExpireTime    int64     `gorm:"type:int(11);size:11;not null;default:0;column:expire_time;comment:过期时间" json:"expire_time,omitempty"`
	Operator      int       `gorm:"type:int(11);size:11;column:operator;comment:更新人/操作人" json:"operator,omitempty"`
	CreateTime    time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;column:create_time;comment:创建时间" json:"create_time"`
	UpdateTime    time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;column:update_time;comment:修改时间" json:"update_time"`
}

func (a *ApplicationAuth) TableName() string {
	return "application_auth"
}

func (a *ApplicationAuth) IdValue() int {
	return a.Id
}

type ApplicationAuthVersion struct {
	Id                int
	ApplicationAuthID int
	NamespaceID       int
	ApplicationAuthVersionConfig
	Operator   int
	CreateTime time.Time
}

type ApplicationAuthVersionConfig struct {
	Config string `json:"config"`
}

func (a *ApplicationAuthVersion) SetVersionId(id int) {
	a.Id = id
}

type ApplicationAuthStat struct {
	ApplicationAuthId int
	VersionID         int
}

// ApplicationAuthRuntime 集群当前版本
type ApplicationAuthRuntime struct {
	Id                int       `json:"id"`
	ClusterId         int       `json:"cluster_id"`
	ApplicationAuthId int       `json:"application_auth_id"`
	NamespaceId       int       `json:"namespace_id"`
	VersionId         int       `json:"version_id"`
	IsOnline          bool      `json:"is_online"`
	Disable           bool      `json:"disable"`
	Operator          int       `json:"operator"`
	CreateTime        time.Time `json:"create_time"`
	UpdateTime        time.Time `json:"update_time"`
}

func (a *ApplicationAuthRuntime) SetRuntimeId(id int) {
	a.Id = id
}

type ApplicationAuthHistoryInfo struct {
	Auth   ApplicationAuth              `json:"auth,omitempty"`
	Config ApplicationAuthVersionConfig `json:"config,omitempty"`
}

type ApplicationAuthHistory struct {
	Id                int
	ApplicationAuthId int
	NamespaceId       int
	OldValue          ApplicationAuthHistoryInfo
	NewValue          ApplicationAuthHistoryInfo
	OptType           history_entry.OptType //1新增 2修改 3删除
	OptTime           time.Time
	Operator          int
}
