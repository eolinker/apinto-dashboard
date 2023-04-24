package dynamic_controller

type DynamicDefine struct {
	Profession string            `json:"profession"`
	Drivers    []*Basic          `json:"drivers"`
	Fields     []*Basic          `json:"fields"`
	Skill      string            `json:"skill"`
	Render     map[string]string `json:"render"`
}

type Render map[string]interface{}

type Basic struct {
	Name  string   `json:"moduleName"`
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
