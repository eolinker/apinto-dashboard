package dto

import (
	"encoding/json"
	"errors"
	"github.com/eolinker/apinto-dashboard/modules/middleware/model"
	"strings"
)

type MiddlewaresInput struct {
	Groups []*model.Middleware `json:"groups"`
}

func (m *MiddlewaresInput) String() string {
	data, _ := json.Marshal(m.Groups)
	return string(data)
}

func (m *MiddlewaresInput) ValidCheck() error {
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
