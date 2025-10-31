package dto

type PluginGroup struct {
	UUID  string `json:"uuid"`
	Name  string `json:"name"`
	Count int    `json:"count,omitempty"`
}
