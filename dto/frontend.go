package dto

type Router struct {
	Name   string            `json:"name"`
	Params map[string]string `json:"params"`
}
