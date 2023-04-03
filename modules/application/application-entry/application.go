package application_entry

import (
	"github.com/eolinker/apinto-dashboard/modules/base/history-entry"
	"time"
)

type Application struct {
	Id          int       `gorm:"type:int(11);size:11;not null;auto_increment;primary_key;column:id;comment:主键ID" json:"id,omitempty"`
	NamespaceId int       `gorm:"type:int(11);size:11;not null;column:namespace;dbUniqueIndex:unique;uniqueIndex:namespace_name;uniqueIndex:namespace_ids;comment:工作空间" json:"namespace_id,omitempty"`
	IdStr       string    `gorm:"size:50;not null;column:id_str;dbUniqueIndex:unique;uniqueIndex:namespace_ids;comment:随机生成的16个长度字符串" json:"id_str,omitempty"`
	Name        string    `gorm:"size:255;not null;column:name;uniqueIndex:namespace_name;comment:应用名称" json:"name,omitempty"`
	Desc        string    `gorm:"size:255;column:desc;comment:描述" json:"desc,omitempty"`
	Operator    int       `gorm:"type:int(11);size:11;column:operator;comment:更新人/操作人" json:"operator,omitempty"`
	CreateTime  time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;column:create_time;comment:创建时间" json:"create_time"`
	UpdateTime  time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;column:update_time;comment:修改时间" json:"update_time"`
}

func (a *Application) TableName() string {
	return "application"
}

func (a *Application) IdValue() int {
	return a.Id
}

type ApplicationVersionConfig struct {
	CustomAttrList []ApplicationCustomAttr `json:"custom_attr_list"`
	ExtraParamList []ApplicationExtraParam `json:"extra_param_list"`
	Apis           []string                `json:"apis"`
}

type ApplicationCustomAttr struct {
	Key   string
	Value string
}

type ApplicationExtraParam struct {
	Key      string
	Value    string
	Position string //header   ["header","body","query"]
	Conflict string //convert  ["origin","convert","error"]
}

type ApplicationVersion struct {
	Id            int
	ApplicationID int
	NamespaceID   int
	ApplicationVersionConfig
	Operator   int
	CreateTime time.Time
}

func (v *ApplicationVersion) SetVersionId(id int) {
	v.Id = id
}

type ApplicationStat struct {
	ApplicationID int
	VersionID     int
}

// ApplicationRuntime 集群当前版本
type ApplicationRuntime struct {
	Id            int       `json:"id"`
	ClusterId     int       `json:"cluster_id"`
	ApplicationId int       `json:"application_id"`
	NamespaceId   int       `json:"namespace_id"`
	VersionId     int       `json:"version_id"`
	IsOnline      bool      `json:"is_online"`
	Disable       bool      `json:"enable"`
	Operator      int       `json:"operator"`
	CreateTime    time.Time `json:"create_time"`
	UpdateTime    time.Time `json:"update_time"`
}

type ApplicationHistoryInfo struct {
	Application              Application              `json:"application"`
	ApplicationVersionConfig ApplicationVersionConfig `json:"config"`
}

type ApplicationHistory struct {
	Id            int
	ApplicationId int
	NamespaceId   int
	OldValue      ApplicationHistoryInfo
	NewValue      ApplicationHistoryInfo
	OptType       history_entry.OptType //1新增 2修改 3删除
	OptTime       time.Time
	Operator      int
}
