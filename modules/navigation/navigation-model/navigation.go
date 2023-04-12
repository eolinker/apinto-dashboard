package navigation_model

type Navigation struct {
	Uuid  string `json:"uuid" yaml:"id"`
	Name  string `json:"-" yaml:"name"`
	Title string `json:"title" yaml:"cname"`
	Icon  string `json:"icon" yaml:"icon"`
}
