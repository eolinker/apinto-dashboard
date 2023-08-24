package dynamic_entry

type DynamicQuote struct {
	Id        int    `gorm:"column:id;type:INT(11);AUTO_INCREMENT;NOT NULL;comment:主键ID"`
	Namespace int    `gorm:"column:namespace;type:INT(11);NOT NULL;comment:工作空间"`
	Source    string `gorm:"column:source;type:VARCHAR(255);NOT NULL;comment:源name"`
	Target    string `gorm:"column:target;type:VARCHAR(255);NOT NULL;comment:依赖name"`
}

func (d *DynamicQuote) TableName() string {
	return "dynamic_quote"
}

func (d *DynamicQuote) IdValue() int {
	return d.Id
}
