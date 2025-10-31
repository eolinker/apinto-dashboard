package pm3

type PluginDefine struct {
	Id          string                  `json:"id,omitempty" yaml:"id"`
	Name        string                  `json:"name,omitempty" yaml:"name"`
	Cname       string                  `json:"cname,omitempty" yaml:"cname"`
	Resume      string                  `json:"resume,omitempty" yaml:"resume"`
	Version     string                  `json:"version,omitempty" yaml:"version"`
	ICon        string                  `json:"icon,omitempty" yaml:"icon"`
	Driver      string                  `json:"driver,omitempty" yaml:"driver"`
	GroupId     string                  `json:"group_id,omitempty" yaml:"group_id"`
	Frontend    []PFrontend             `json:"frontend,omitempty" yaml:"frontend"`
	Navigations []NavigationItem        `json:"navigations,omitempty" yaml:"navigations"`
	Define      interface{}             `json:"define,omitempty" yaml:"define"`
	Access      map[string][]AccessItem `json:"access,omitempty" yaml:"access"`
	Depend      []DependItem            `json:"depend,omitempty" yaml:"depend"`
}

type DependItem struct {
	Id    string `json:"id,omitempty" yaml:"id"`
	Cname string `json:"cname,omitempty" yaml:"cname"`
	Link  string `json:"link,omitempty" yaml:"link"`
}
type NavigationItem struct {
	Navigation string       `json:"navigation,omitempty" yaml:"navigation"`
	Router     string       `json:"router,omitempty" yaml:"router"`
	Name       string       `json:"name,omitempty" yaml:"name"`
	Cname      string       `json:"cname,omitempty" yaml:"cname"`
	Access     []AccessItem `json:"access,omitempty" yaml:"access"`
}

type AccessItem struct {
	Name   string   `yaml:"name" json:"name,omitempty"`
	Cname  string   `yaml:"cname" json:"cname,omitempty"`
	Depend []string `yaml:"depend" json:"depend,omitempty"`
}
