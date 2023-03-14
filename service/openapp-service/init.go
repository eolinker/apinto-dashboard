package openapp_service

import "github.com/eolinker/eosc/common/bean"

func init() {
	externalAPP := newExternalApplicationService()

	bean.Injection(&externalAPP)
}
