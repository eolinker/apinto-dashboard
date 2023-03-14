package variable_model

import (
	variable_entry2 "github.com/eolinker/apinto-dashboard/modules/variable/variable-entry"
)

type GlobalVariable struct {
	*variable_entry2.Variables
}

type GlobalVariableListItem struct {
	*variable_entry2.Variables
	Status      int
	OperatorStr string
}

type GlobalVariableDetails struct {
	*variable_entry2.ClusterVariable
	Status      int
	ClusterName string
	Environment string
}
