package dto

type CommonGroupOut struct {
	UUID     string            `json:"uuid"`
	Name     string            `json:"name"`
	Children []*CommonGroupOut `json:"children"` //子目录
	IsDelete bool              `json:"is_delete"`
}

type CommonGroupRootOut struct {
	UUID     string            `json:"uuid"`
	Name     string            `json:"name"`
	Groups   []*CommonGroupOut `json:"groups"`
	IsDelete bool              `json:"is_delete"`
}

type CommonGroupInput struct {
	Name       string `json:"name"`
	UUID       string `json:"uuid"`
	ParentUUID string `json:"parent_uuid"`
}

type CommonGroupApi struct {
	Name      string   `json:"name"`
	UUID      string   `json:"uuid"`
	Methods   []string `json:"methods"`
	GroupUUID string   `json:"group_uuid"`
}

type CommGroupSortInput struct {
	Root  string   `json:"root"`  //操作目标id，根目录为空
	Items []string `json:"items"` //操作之后，root下级的目录uuid数组
}
