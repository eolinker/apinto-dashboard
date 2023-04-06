package dto

import (
	"encoding/json"
	"errors"
	"strings"
)

type Middlewares struct {
	Groups []*Middleware `json:"groups"`
}

func (m *Middlewares) String() string {
	data, _ := json.Marshal(m.Groups)
	return string(data)
}

func (m *Middlewares) ValidCheck() error {
	prefixMap := map[string]struct{}{}
	for _, g := range m.Groups {
		prefix := strings.TrimSuffix(strings.TrimSpace(g.Prefix), "/")
		if _, ok := prefixMap[prefix]; ok {
			return errors.New("repeat prefix: " + g.Prefix)
		}
		prefixMap[prefix] = struct{}{}
	}
	return nil
}

type Middleware struct {
	Prefix      string   `json:"prefix"`
	Middlewares []string `json:"middlewares"`
}
