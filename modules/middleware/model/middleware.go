package model

type Middleware struct {
	Prefix      string   `json:"prefix"`
	Middlewares []string `json:"middlewares"`
}
