package plugin_entry

import "time"

type Plugin struct {
	Id          int       `gorm:"type:int(11);size:11;not null;auto_increment;primary_key;column:id;comment:主键ID" json:"id,omitempty"`
	NamespaceId int       `gorm:"type:int(11);size:11;not null;column:namespace;dbUniqueIndex:namespace_name_index;uniqueIndex:namespace_name_index;comment:工作空间" json:"namespace_id,omitempty"`
	Name        string    `gorm:"size:36;not null;column:name;dbUniqueIndex:namespace_name_index;uniqueIndex:namespace_name_index;comment:插件名称" json:"name,omitempty"`
	Extended    string    `gorm:"size:255;column:extended;comment:扩展ID" json:"extended,omitempty"`
	Desc        string    `gorm:"size:255;column:desc;comment:描述" json:"desc,omitempty"`
	Schema      string    `gorm:"column:schema;comment:jsonSchema" json:"schema,omitempty"`
	Type        int       `gorm:"column:type;comment:插件类型内置1 自建2" json:"type,omitempty"`
	Rely        int       `gorm:"column:rely;comment:依赖的插件ID" json:"rely,omitempty"`
	Sort        int       `gorm:"column:sort;comment:排序字段" json:"sort,omitempty"`
	Operator    int       `gorm:"type:int(11);size:11;column:operator;comment:更新人/操作人" json:"operator,omitempty"`
	CreateTime  time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;column:create_time;comment:创建时间" json:"create_time"`
	UpdateTime  time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;column:update_time;comment:修改时间" json:"update_time"`
}

func (*Plugin) TableName() string {
	return "plugin"
}

func (p *Plugin) IdValue() int {
	return p.Id
}
