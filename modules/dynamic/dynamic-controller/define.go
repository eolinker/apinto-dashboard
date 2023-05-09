package dynamic_controller

type OptionTitle struct {
	Field string `json:"field"`
	Title string `json:"title"`
}
type FilterOptionConfig struct {
	Name   string        `json:"name"`
	Title  string        `json:"title"`
	Titles []OptionTitle `json:"titles"`
}
type DynamicDefine struct {
	Profession    string              `json:"profession"`
	Drivers       []*Basic            `json:"drivers"`
	Fields        []*Basic            `json:"fields"`
	Skill         string              `json:"skill"`
	FilterOptions *FilterOptionConfig `json:"options"`
	Render        map[string]string   `json:"render"`
}

type Render map[string]interface{}

type Basic struct {
	Name  string   `json:"name"`
	Title string   `json:"title"`
	Attr  string   `json:"attr,omitempty"`
	Enum  []string `json:"enum,omitempty"`
}

var defaultFields = []*Basic{
	{
		Name:  "updater",
		Title: "更新者",
	},
	{
		Name:  "update_time",
		Title: "更新时间",
	},
}
