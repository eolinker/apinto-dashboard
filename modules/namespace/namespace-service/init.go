package namespace_service

import "github.com/eolinker/eosc/common/bean"

func init() {
	namespace := newNamespaceService()

	bean.Injection(&namespace)
}
