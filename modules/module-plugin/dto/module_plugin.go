package dto

type PluginListItem struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Cname   string `json:"cname"`
	Resume  string `json:"resume"`
	ICon    string `json:"icon"`
	Enable  bool   `json:"enable"`
	IsInner bool   `json:"is_inner"`
}

type PluginGroup struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type PluginInfo struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Cname     string `json:"cname"`
	Resume    string `json:"resume"`
	Icon      string `json:"icon"`
	Enable    bool   `json:"enable"`
	Uninstall bool   `json:"uninstall"`
}

type PluginEnableInfo struct {
	Name       string         `json:"name"`
	Server     string         `json:"server"`
	Header     []ExtendParams `json:"header"`
	Query      []ExtendParams `json:"query"`
	Initialize []ExtendParams `json:"initialize"`
}

type PluginEnableRender struct {
	Internet   bool                 `json:"internet"`
	Headers    []ExtendParamsRender `json:"headers"`
	Querys     []ExtendParamsRender `json:"querys"`
	Initialize []ExtendParamsRender `json:"initialize"`
}

type ExtendParams struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ExtendParamsRender struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	Title       string `json:"title"`
	Placeholder string `json:"placeholder"`
	Desc        string `json:"desc"`
}
