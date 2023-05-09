package random_service

import "github.com/eolinker/eosc/common/bean"

func init() {
	random := newRandomService()

	bean.Injection(&random)
}
