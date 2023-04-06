package entry

type System struct {
	Id    int    `gorm:"column:id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (s *System) TableName() string {
	return "system_info"
}

func (s *System) IdValue() int {
	return s.Id
}
