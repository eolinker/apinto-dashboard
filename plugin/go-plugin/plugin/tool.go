package plugin

import (
	"github.com/eolinker/apinto-dashboard/plugin/go-plugin/shared"
	"github.com/eolinker/apinto-dashboard/pm3"
)

var (
	_ shared.Middleware = (*simpleMiddleware)(nil)
)

type simpleMiddleware struct {
	checkHandler    func(info pm3.ApiInfo) bool
	requestHandler  shared.MiddlewareHandFunc
	hasRequest      bool
	responseHandler shared.MiddlewareHandFunc
	hasResponse     bool
}

func (s *simpleMiddleware) RequestHandler() (shared.MiddlewareHandFunc, bool) {

	return s.requestHandler, s.hasRequest
}

func (s *simpleMiddleware) ResponseHandler() (shared.MiddlewareHandFunc, bool) {
	return s.responseHandler, s.hasResponse
}

func (s *simpleMiddleware) Check(api pm3.ApiInfo) bool {
	return s.checkHandler(api)
}

type SetFunc func(builder *simpleMiddleware)

func NewMiddleware(checkH func(info pm3.ApiInfo) bool, setFuncs ...SetFunc) shared.Middleware {
	m := &simpleMiddleware{
		checkHandler: checkH,
	}
	for _, setFunc := range setFuncs {
		setFunc(m)
	}
	return m

}

func ProcessRequestBy(f shared.MiddlewareHandFunc) SetFunc {
	return func(m *simpleMiddleware) {
		m.requestHandler = f
		m.hasRequest = f != nil
	}
}
func ProcessResponseBy(f shared.MiddlewareHandFunc) SetFunc {
	return func(m *simpleMiddleware) {
		m.responseHandler = f
		m.hasResponse = f != nil
	}
}
