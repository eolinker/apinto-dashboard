package variable_service

import "github.com/eolinker/eosc/common/bean"

func init() {
	globalVariable := newGlobalVariableService()

	bean.Injection(&globalVariable)

	clusterVariable := newClusterVariableService()

	bean.Injection(&clusterVariable)
}
