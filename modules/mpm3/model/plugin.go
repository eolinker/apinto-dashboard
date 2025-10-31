package model

type Plugin struct {
	UUID    string
	Name    string
	CName   string
	Resume  string
	ICon    string
	Enable  bool
	IsInner bool

	Group string
}
type PluginInfo struct {
	Plugin
	Uninstall  bool
	CanDisable bool
}
type PluginDetail struct {
}
