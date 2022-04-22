package apinto_dashboard

type ActivityLogAddHandler func(user, content, operation, object string, args []*Arg) error
type ActivityLogGetHandler func(offset, limit int) ([]*LogEntity, int64, error)

type LogEntity struct {
	Time      string `json:"time"`
	User      string `json:"user"`
	Operation string `json:"operation"`
	Object    string `json:"object"`
	Content   string `json:"content"`
	Args      []Arg  `json:"args"`
}

type Arg struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

func AddActivityLog(user, content, operation) {

}
