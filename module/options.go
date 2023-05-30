package apinto_module

type IFilterOptionHandlerSupport interface {
	FilterOptionHandler() []IFilterOptionHandler
}
type FilterOptionConfig struct {
	Title string

	Titles []OptionTitle
	Key    string
}
type IFilterValue struct {
}
type IFilterOptionHandler interface {
	Name() string
	Config() FilterOptionConfig
	GetOptions(namespaceId int, keyword, groupUUID string, pageNum, pageSize int) ([]any, int)
	Labels(namespaceId int, values ...string) []string
	Label(namespaceId int, value string) string
}

type OptionTitle struct {
	Field string `json:"field"`
	Title string `json:"title"`
}

type IFilterOptionHandlerManager interface {
	ResetFilterOptionHandlers(handlers map[string]IFilterOptionHandler)
}
