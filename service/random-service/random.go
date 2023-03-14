package random_service

import "github.com/eolinker/apinto-dashboard/common"

type IRandomService interface {
	RandomStr(ruleName string) string
}
type randomService struct {
}

func newRandomService() IRandomService {
	return &randomService{}
}

func (randomService) RandomStr(ruleName string) string {
	// todo 这里需要换成通过注册取代来实现
	switch ruleName {
	case "application":
		return common.RandStr(16)
	case "external-app":
		return common.RandStr(16)
	case "external-app-token": //外部应用长度32的token
		return common.RandStr(32)
	case "password":
		return "12345678"
		return common.RandStr(8)
	default:
		return ""
	}
}
