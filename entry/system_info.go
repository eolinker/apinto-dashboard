package entry

// SystemInfo 系统信息表
type SystemInfo struct {
	Id    int    `gorm:"type:int(11);size:11;not null;auto_increment;primary_key;column:id;comment:主键ID"`
	Key   string `gorm:"size:20;not null;column:key;dbUniqueIndex:unique_key;uniqueIndex:unique_key;comment:健"`
	Value []byte `gorm:"type:text;not null;column:value;comment:值"`
}

func (*SystemInfo) TableName() string {
	return "system_info"
}

func (c *SystemInfo) IdValue() int {
	return c.Id
}
