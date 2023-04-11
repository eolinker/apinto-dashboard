package navigation_entry

type Navigation struct {
	Id       int    `gorm:"column:id"`
	Uuid     string `gorm:"column:uuid"`
	Title    string `gorm:"column:title"`
	Icon     string `gorm:"column:icon"`
	IconType string `gorm:"column:icon_type"`
	Sort     int    `gorm:"column:sort"`
	Module   string `gorm:"column:module"`
}

func (n *Navigation) TableName() string {
	return "navigation"
}

func (n *Navigation) IdValue() int {
	return n.Id
}
