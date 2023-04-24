package local

type Define struct {
	Middleware []MiddlewareConfig `json:"middleware"`
	Router     RouterConfig       `json:"router"`
	Provider   []ProviderConfig   `json:"provider"`
}
type MiddlewareConfig struct {
	Name string   `json:"name"`
	Path string   `json:"path"`
	Life string   `json:"life"`
	Rule []string `json:"rule"`
}
type RouterConfig struct {
	Home     string                          `json:"home"`
	Html     []PathConfig                    `json:"html"`
	Frontend []string                        `json:"frontend"`
	Api      map[string]map[string]Attribute `json:"api"`
	OpenApi  map[string]map[string]Attribute `json:"openapi"`
	Provider []ProviderConfig
}
type PathConfig struct {
	Label []string `json:"label"`
	Path  string   `json:"path"`
}
type Attribute struct {
	Label []string `json:"label"`
}

type ProviderConfig struct {
	Name string `json:"name"`
	Path string `json:"path"`
}
