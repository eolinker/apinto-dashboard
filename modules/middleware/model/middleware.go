package model

type Middleware struct {
	Group       []*MiddlewareGroup `json:"group"`
	Middlewares []*MiddlewareInfo  `json:"middlewares"`
}

type MiddlewareGroup struct {
	Prefix      string   `json:"prefix"`
	Middlewares []string `json:"middlewares"`
}

type MiddlewareInfo struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}
