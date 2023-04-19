package model

type Navigation struct {
	Title string `json:"title"`
	Icon  string `json:"icon"`
	//IconType string    `json:"icon_type"`
	Modules []*Module `json:"modules"`
	Default string    `json:"default"`
}

type Module struct {
	Name  string `json:"name"`
	Title string `json:"title"`
	Type  string `json:"type"`
	Path  string `json:"path"`
}
