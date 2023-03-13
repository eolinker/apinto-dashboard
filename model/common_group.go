package model

import "github.com/eolinker/apinto-dashboard/entry"

type CommonGroupRoot struct {
	UUID        string
	Name        string
	CommonGroup []*CommonGroup
}

type CommonGroup struct {
	Group    *entry.CommonGroup
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
