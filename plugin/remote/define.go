package remote

type Define struct {
	Internet   bool                 `json:"internet" yaml:"internet"`
	Server     string               `json:"server" yaml:"server"`
	Path       string               `json:"path" yaml:"path"`
	Headers    []ExtendParamsRender `json:"headers,omitempty" yaml:"headers,omitempty"`       //local
	Querys     []ExtendParamsRender `json:"querys,omitempty" yaml:"querys,omitempty"`         //remote local
	Initialize []ExtendParamsRender `json:"initialize,omitempty" yaml:"initialize,omitempty"` //remote local
}

type ExtendParamsRender struct {
	Name        string `json:"name" yaml:"name"`
	Value       string `json:"value" yaml:"value"`
	Title       string `json:"title" yaml:"title"`
	Type        string `json:"type" yaml:"type"`
	Placeholder string `json:"placeholder" yaml:"placeholder"`
	Desc        string `json:"desc" yaml:"desc"`
}
