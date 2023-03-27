package plugin_dto

type PluginListItem struct {
	Name       string `json:"name,omitempty"`
	Extended   string `json:"extended,omitempty"`
	Desc       string `json:"desc,omitempty"`
	UpdateTime string `json:"update_time,omitempty"`
	Operator   string `json:"operator,omitempty"`
	IsDelete   bool   `json:"is_delete"`
	IsBuilt    bool   `json:"is_built"`
	Config     string `json:"config,omitempty"`
}

type PluginItem struct {
	Name     string `json:"name"`
	Extended string `json:"extended"`
	Desc     string `json:"desc"`
	Rely     string `json:"rely"`
}

type PluginInput struct {
	Name     string `json:"name"`
	Extended string `json:"extended"`
	Rely     string `json:"rely"`
	Desc     string `json:"desc"`
}

type PluginSort struct {
	Names []string `json:"names"`
}

type PluginEnum struct {
	Name   string `json:"name"`
	Config string `json:"config"`
}
