package plugin_template_entry

import "time"

type PluginTemplate struct {
	Id          int       `gorm:"type:int(11);size:11;not null;auto_increment;primary_key;column:id;comment:主键ID" json:"id,omitempty"`
	NamespaceId int       `gorm:"type:int(11);size:11;not null;column:namespace;comment:工作空间" json:"namespace_id,omitempty"`
	UUID        string    `gorm:"size:36;not null;column:uuid;dbUniqueIndex:uuid;comment:UUID" json:"uuid,omitempty"`
	Name        string    `gorm:"size:255;column:name;comment:插件模板名称" json:"name,omitempty"`
	Desc        string    `gorm:"size:255;column:desc;comment:描述" json:"desc,omitempty"`
	Operator    int       `gorm:"type:int(11);size:11;column:operator;comment:更新人/操作人" json:"operator,omitempty"`
	CreateTime  time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;column:create_time;comment:创建时间" json:"create_time"`
	UpdateTime  time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;column:update_time;comment:修改时间" json:"update_time"`
}

func (*PluginTemplate) TableName() string {
	return "plugin_template"
}

func (p *PluginTemplate) IdValue() int {
	return p.Id
}
