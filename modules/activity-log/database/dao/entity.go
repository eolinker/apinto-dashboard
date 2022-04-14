package dao

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
