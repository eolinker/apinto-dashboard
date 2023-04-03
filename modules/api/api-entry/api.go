package api_entry

import (
	"time"
)

type API struct {
	Id               int       `gorm:"type:int(11);size:11;not null;auto_increment;primary_key;column:id;comment:主键ID" json:"id,omitempty"`
	NamespaceId      int       `gorm:"type:int(11);size:11;not null;column:namespace;uniqueIndex:namespace_uuid;comment:工作空间" json:"namespace_id,omitempty"`
	UUID             string    `gorm:"size:36;not null;column:uuid;dbUniqueIndex:namespace_uuid;uniqueIndex:namespace_uuid;comment:UUID" json:"uuid,omitempty"`
	GroupUUID        string    `gorm:"size:36;column:group_uuid;comment:api所在分组的UUID;index:group_uuid" json:"group_uuid,omitempty"`
	Name             string    `gorm:"size:255;column:name;comment:api名称" json:"name,omitempty"`
	RequestPath      string    `gorm:"size:255;request_path;comment:api请求路径" json:"request_path,omitempty"`
	RequestPathLabel string    `gorm:"size:255;request_path_label;comment:api请求路径Label" json:"request_path_label,omitempty"`
	SourceType       string    `gorm:"size:255;source_type;comment:来源类型" json:"source_type,omitempty"`
	SourceID         int       `gorm:"type:int(11);size:11;column:source_id;comment:来源id,用于关联外部应用" json:"source_id,omitempty"`
	SourceLabel      string    `gorm:"size:255;source_label;comment:来源标签" json:"source_label,omitempty"`
	Desc             string    `gorm:"size:255;column:desc;comment:描述" json:"desc,omitempty"`
	Operator         int       `gorm:"type:int(11);size:11;column:operator;comment:更新人/操作人" json:"operator,omitempty"`
	CreateTime       time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;column:create_time;comment:创建时间" json:"create_time"`
	UpdateTime       time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;column:update_time;comment:修改时间" json:"update_time"`
}

// /user/{uuid}
// path:/user/{rest}
// PathLabel:/user/{uuid}

func (*API) TableName() string {
	return "api"
}

func (s *API) IdValue() int {
	return s.Id
}

type APISource struct {
	SourceType  string
	SourceID    int
	SourceLabel string
}

type APISourceList []*APISource

func (a APISourceList) Len() int {
	return len(a)
}

func (a APISourceList) Less(i, j int) bool {
	//按自建，导入，同步顺序排，同步按应用id升序，id相同按字母升序
	if a[i].SourceID < a[j].SourceID {
		return true
	} else if a[i].SourceLabel < a[j].SourceLabel {
		return true
	}
	return a[i].SourceType == "self-build"
}

func (a APISourceList) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
