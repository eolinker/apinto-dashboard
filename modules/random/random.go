package random

type IRandomService interface {
	RandomStr(ruleName string) string
}
