package variable_model

import (
	"github.com/eolinker/apinto-dashboard/modules/variable/variable-entry"
)

type GlobalVariable struct {
	*variable_entry.Variables
}

type GlobalVariableListItem struct {
	*variable_entry.Variables
	Status      int
	OperatorStr string
}

func (g *GlobalVariableListItem) UserId() int {
	return g.Variables.Operator
}

func (g *GlobalVariableListItem) Set(name string) {
	g.OperatorStr = name
}

type GlobalVariableDetails struct {
	*variable_entry.ClusterVariable
	Status      int
	ClusterName string
	Environment string
}
