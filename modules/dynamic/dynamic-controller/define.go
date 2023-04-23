package dynamic_controller

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
