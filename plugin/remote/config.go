package remote

type Config struct {
	Server     string            `json:"server"`
	Header     map[string]string `json:"header"`
	Query      map[string]string `json:"query"`
	Initialize map[string]string `json:"initialize"`
}

type DtoOpenModeParam struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Type  string `json:"type,omitempty"`
}
