package local

type Config struct {
	Server     string            `json:"server"`
	Header     map[string]string `json:"header"`
	Query      map[string]string `json:"query"`
	Initialize map[string]string `json:"initialize"`
}
