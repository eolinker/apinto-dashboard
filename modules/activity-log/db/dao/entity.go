package dao

type LogEntity struct {
	Time    string `json:"time"`
	User    string `json:"user"`
	Content string `json:"content"`
	Args    []*Arg `json:"args"`
}

type Arg struct {
	key   string
	value interface{}
}
