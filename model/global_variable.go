package model

import "github.com/eolinker/apinto-dashboard/entry"

type GlobalVariable struct {
	*entry.Variables
}

type GlobalVariableListItem struct {
	*entry.Variables
	Status      int
	OperatorStr string
}

type GlobalVariableDetails struct {
	*entry.ClusterVariable
	Status      int
	ClusterName string
	Environment string
}
