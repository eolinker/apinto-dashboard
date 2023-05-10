package dynamic_entry

type DynamicQuote struct {
	Id        int    `gorm:"column:id"`
	Namespace int    `gorm:"column:namespace"`
	Source    string `gorm:"column:source"`
	Target    string `gorm:"column:target"`
}

func (d *DynamicQuote) TableName() string {
	return "dynamic_quote"
}

func (d *DynamicQuote) IdValue() int {
	return d.Id
}
