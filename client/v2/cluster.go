package v2

type Cluster struct {
	Cluster string  `yaml:"cluster"`
	Nodes   []*Node `yaml:"nodes"`
}

type Node struct {
	Id     string   `json:"id"`
	Name   string   `json:"name"`
	Peer   []string `yaml:"peer"`
	Admin  []string `json:"admin"`
	Server []string `json:"server"`
	Leader bool     `json:"leader"`
}
