package dynamic_model

type DynamicListItem struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Title      string `json:"title"`
	Driver     string `json:"driver"`
	Updater    string `json:"updater"`
	UpdateTime string `json:"update_time"`
}

// DynamicDriver 动态模块驱动信息
type DynamicDriver struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}

type DynamicField struct {
	Name  string   `json:"name"`
	Title string   `json:"title"`
	Attr  string   `json:"attr,omitempty"`
	Enum  []string `json:"enum,omitempty"`
}

type DynamicInfo struct {
	*DynamicBasicInfo
	Append map[string]interface{}
}

type DynamicBasicInfo struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Online      bool   `json:"-"`
}

// DynamicCluster 动态模块集群信息
type DynamicCluster struct {
	Name       string `json:"name"`
	Title      string `json:"title"`
	Status     string `json:"status"`
	Updater    string `json:"updater"`
	UpdateTime string `json:"update_time"`
}
