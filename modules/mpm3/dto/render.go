package dto

type PluginEnableRender struct {
	Internet bool `json:"internet"`
	//Invisible  bool                 `json:"invisible"`
	NameConflict bool                 `json:"name_conflict"`
	Headers      []ExtendParamsRender `json:"headers"`
	Querys       []ExtendParamsRender `json:"querys"`
	Initialize   []ExtendParamsRender `json:"initialize"`
}

type ExtendParamsRender struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	Title       string `json:"title"`
	Type        string `json:"type,omitempty"`
	Placeholder string `json:"placeholder"`
	Desc        string `json:"desc"`
}
