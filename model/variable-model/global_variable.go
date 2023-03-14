package variable_model

import (
	"github.com/eolinker/apinto-dashboard/entry/variable-entry"
)

type GlobalVariable struct {
	*variable_entry.Variables
}

type GlobalVariableListItem struct {
	*variable_entry.Variables
	Status      int
	OperatorStr string
}

type GlobalVariableDetails struct {
	*variable_entry.ClusterVariable
	Status      int
	ClusterName string
	Environment string
}
