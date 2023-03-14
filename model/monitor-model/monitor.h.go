package monitor_model

type MonWhereItem struct {
	Key       string
	Operation string // 表达式，默认为 =，多个为 in，可以用其他
	Values    []string
}
type MonSortBy struct {
	Key  string
	Desc bool
}
type MonStatisticsValue struct {
	MonCommonData
}
type MonTrendFilter struct {
	Name string
	MonWhereItem
}

type MonTrendValues struct {
	Data   []string
	Names  []string
	Values [][]interface{}
}
