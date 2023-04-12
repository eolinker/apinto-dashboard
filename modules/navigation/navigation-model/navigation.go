package navigation_model

type Navigation struct {
	Uuid  string `json:"uuid" yaml:"id"`
	Title string `json:"title" yaml:"name"`
	Icon  string `json:"icon" yaml:"icon"`
}
