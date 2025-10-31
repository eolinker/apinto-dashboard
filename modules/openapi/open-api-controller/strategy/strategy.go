package strategy

type IStrategyValidator interface {
	Validate(data []byte) ([]byte, error)
}

type RewriteResponse struct {
	Body        string `json:"body"`
	Charset     string `json:"charset"`
	ContentType string `json:"content_type"`
	StatusCode  int    `json:"status_code"`
}

type Filter struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

var (
	validators = map[string]IStrategyValidator{
		"visit":    &VisitStrategyValidator{},
		"limiting": &LimitStrategyValidator{},
	}
)

func GetValidator(name string) (IStrategyValidator, bool) {
	v, ok := validators[name]
	return v, ok
}
