package model

type RoleListItem struct {
	ID             string
	Title          string
	UserNum        int
	OperateDisable bool
	Type           int
}

type RoleInfo struct {
	Title  string
	Desc   string
	Access []string
}

type RoleOptionItem struct {
	ID             string
	Title          string
	OperateDisable bool
}
