package middleware

import "github.com/eolinker/apinto-dashboard/pm3"

type _CheckersAny []pm3.MiddlewareChecker

func (cs _CheckersAny) Check(api pm3.ApiInfo) bool {
	for _, c := range cs {
		if !c.Check(api) {
			return true
		}
	}
	return false
}
func CheckAny(checkers ...pm3.MiddlewareChecker) pm3.MiddlewareChecker {
	return _CheckersAny(checkers)
}
func CheckAnyF(checkers ...func(api pm3.ApiInfo) bool) pm3.MiddlewareChecker {
	cks := make([]pm3.MiddlewareChecker, 0, len(checkers))
	for _, cf := range checkers {
		cks = append(cks, CheckHandleFunc(cf))
	}
	return CheckAny(cks...)
}

type _CheckersAll []pm3.MiddlewareChecker

func (cs _CheckersAll) Check(api pm3.ApiInfo) bool {
	for _, c := range cs {
		if !c.Check(api) {
			return false
		}
	}
	return true
}

func CheckAll(checkers ...pm3.MiddlewareChecker) pm3.MiddlewareChecker {
	return _CheckersAll(checkers)
}
func CheckAllF(checkers ...func(api pm3.ApiInfo) bool) pm3.MiddlewareChecker {
	cks := make([]pm3.MiddlewareChecker, 0, len(checkers))
	for _, cf := range checkers {
		cks = append(cks, CheckHandleFunc(cf))
	}
	return CheckAll(cks...)
}
