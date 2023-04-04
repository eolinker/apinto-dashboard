package model

type MiddlewareGroup struct {
	Uuid   string `json:"uuid"`
	Prefix string `json:"prefix"`
}

type MiddlewareGroups []*MiddlewareGroup

func (m MiddlewareGroups) Len() int {
	return len(m)
}

func (m MiddlewareGroups) Less(i, j int) bool {
	return m[i].Prefix > m[j].Prefix
}

func (m MiddlewareGroups) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

type MiddlewareGroupInfo struct {
	Middlewares []string `json:"middlewares"`
	MiddlewareGroup
	All []*MiddlewareInfo `json:"all"`
}

type MiddlewareInfo struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}
