package entry

import "time"

// UserRole 用户关联角色表
type UserRole struct {
	ID         int       `gorm:"type:int(11);size:11;not null;auto_increment;primary_key;column:id;comment:主键ID" json:"id"`
	UserID     int       `gorm:"type:int(11);size:11;not null;column:user_id;dbUniqueIndex:unique;index:user_id;comment:用户id" json:"user_id"`
	RoleID     int       `gorm:"type:int(11);size:11;not null;column:role_id;dbUniqueIndex:unique;index:role_id;comment:角色ID" json:"role_id"`
	Module     string    `gorm:"size:512;not null;column:module;comment:模块目录（多级模块用‘/’切割）" json:"module"`
	CreateTime time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;column:create_time;comment:创建时间" json:"create_time"`
}

func (r *UserRole) IdValue() int {
	return r.RoleID
}

func (*UserRole) TableName() string {
	return "user_role"
}
