package role_entry

import "time"

// RoleAccess 角色关联权限表
type RoleAccess struct {
	RoleID     int       `gorm:"type:int(11);size:11;not null;primary_key;column:role_id;comment:角色ID" json:"role_id"`
	AccessIDs  string    `gorm:"size:255;not null;column:access_ids;comment:权限ID" json:"access_ids"`
	CreateTime time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;column:create_time;comment:创建时间" json:"create_time"`
}

func (r *RoleAccess) IdValue() int {
	return r.RoleID
}

func (*RoleAccess) TableName() string {
	return "role_access"
}
