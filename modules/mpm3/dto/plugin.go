package dto

type PluginEnableInfo struct {
	Name       string         `json:"name"`
	Header     []ExtendParams `json:"header"`
	Query      []ExtendParams `json:"query"`
	Initialize []ExtendParams `json:"initialize"`
}

type ExtendParams struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type PluginInfo struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Cname      string `json:"cname"`
	Resume     string `json:"resume"`
	Icon       string `json:"icon"`
	Enable     bool   `json:"enable"`
	CanDisable bool   `json:"can_disable"`
	Uninstall  bool   `json:"uninstall"`
}

type PluginListItem struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Cname   string `json:"cname"`
	Resume  string `json:"resume"`
	ICon    string `json:"icon"`
	Enable  bool   `json:"enable"`
	IsInner bool   `json:"is_inner"`
}
