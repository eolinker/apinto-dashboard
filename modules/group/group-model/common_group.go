package group_model

import (
	"github.com/eolinker/apinto-dashboard/modules/group/group-entry"
)

type CommonGroupRoot struct {
	UUID        string
	Name        string
	CommonGroup []*CommonGroup
}

type CommonGroup struct {
	Group    *group_entry.CommonGroup
	Subgroup []*CommonGroup
}

type CommonGroupApi struct {
	Path      string
	PathLabel string
	Name      string
	UUID      string
	Methods   []string
	GroupUUID string
}
