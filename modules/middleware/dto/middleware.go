package dto

type MiddlewareGroup struct {
	ID          string   `json:"id"`
	Prefix      string   `json:"prefix"`
	Middlewares []string `json:"middlewares"`
}
