package entry

type EnumType int32

const (
	VariableType EnumType = 1 //环境变量
)

type Enum struct {
	Id    int      `gorm:"primary_key;column:id"`
	Name  string   `gorm:"column:name"`
	Value string   `gorm:"column:value"`
	Type  EnumType `gorm:"column:type"`
	Sort  int      `gorm:"column:sort"`
}

func (*Enum) TableName() string {
	return "enum"
}

func (e *Enum) IdValue() int {
	return e.Id
}
