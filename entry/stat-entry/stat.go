package stat_entry

// Stat 最新版本关联表
type Stat struct {
	ID      int    `gorm:"type:int(11);size:11;not null;auto_increment;primary_key;column:id;comment:主键ID"`
	Kind    string `gorm:"size:20;not null;column:kind;dbUniqueIndex:kind_tag;uniqueIndex:kind_tag;comment:根据kind区分是哪个表的ID"`
	Tag     int    `gorm:"type:int(11);size:11;not null;column:target;dbUniqueIndex:kind_tag;uniqueIndex:kind_tag"`
	Version int    `gorm:"type:int(11);size:11;not null;column:version;comment:版本ID"`
}

func (s *Stat) IdValue() int {
	return s.ID
}

func (*Stat) TableName() string {
	return "stat"
}
